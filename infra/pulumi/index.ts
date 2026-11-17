// Hosting for shell-explain.ninegene.com: a private S3 bucket behind CloudFront (origin
// access control), an ACM certificate, DNS records in the existing ninegene.com zone and
// the GitHub OIDC deploy role.
//
// Adapted from ninegene/ninegene.com (infra/pulumi/index.ts). Differences:
//   - The Route 53 zone is owned (and protected) by the ninegene.com stack. This program
//     only looks it up and adds its own records, so it never changes the zone itself.
//   - The state bucket, KMS key and GitHub OIDC provider come from that repository's
//     infra/bootstrap/bootstrap.sh and are shared.
//   - One host name, no www redirect.
//
// Regions: the AWS account is a project account from AWS's new sign-up with advanced
// features active. Its managed policies allow only global services (CloudFront, IAM,
// Route 53, ACM, STS) in us-east-1, so Regional resources (the site bucket) live in the
// project's Region (`region` config, us-east-2). The CloudFront certificate must be
// requested in us-east-1 and uses its own provider.
// See infra/README.md.
import * as fs from 'node:fs';
import * as path from 'node:path';
import * as aws from '@pulumi/aws';
import * as pulumi from '@pulumi/pulumi';

const config = new pulumi.Config();
const zoneName = config.get('zoneName') ?? 'ninegene.com';
const domain = config.get('domain') ?? `shell-explain.${zoneName}`;
const githubRepo = config.get('githubRepo') ?? 'ninegene/shell-explain';
const githubEnvironment = config.get('githubEnvironment') ?? 'production';
// The start of GitHub's OIDC `sub` claim, which the deploy role's trust policy must match
// exactly. This repository uses immutable subjects (checked 2026-11-16 with
// `gh api repos/ninegene/shell-explain/actions/oidc/customization/sub`), so the stack
// sets `githubSubjectPrefix` to the `sub_claim_prefix` that command returns. The default
// below is the name-based subject, used only if the repository switches back.
const githubSubjectPrefix =
  config.get('githubSubjectPrefix') ?? `repo:${githubRepo}`;
const siteBucketName = config.get('siteBucketName') ?? 'shell-explain-site';
const region = config.get('region') ?? 'us-east-2';

const tags = {
  project: 'shell-explain',
  environment: 'production',
  'managed-by': 'pulumi',
};
const provider = new aws.Provider('aws', {
  region,
  defaultTags: { tags },
});
const opts: pulumi.CustomResourceOptions = { provider };
// CloudFront only accepts ACM certificates from us-east-1.
const usEast1 = new aws.Provider('aws-us-east-1', {
  region: 'us-east-1',
  defaultTags: { tags },
});

const account = aws.getCallerIdentityOutput({}, { provider });
// Created once by ninegene.com's infra/bootstrap/bootstrap.sh, not by this program.
const oidcProviderArn = pulumi.interpolate`arn:aws:iam::${account.accountId}:oidc-provider/token.actions.githubusercontent.com`;

// --- DNS: records in the existing zone ----------------------------------------------
// The zone belongs to the ninegene.com stack; this only reads it.
const zone = aws.route53.getZoneOutput({ name: zoneName }, { provider });

// --- Site bucket: private, versioned, readable only through CloudFront --------------
const bucket = new aws.s3.Bucket(
  'site',
  { bucket: siteBucketName },
  { provider, protect: true },
);
export const siteBucket = bucket.bucket;
// S3 rejects bucket-configuration calls made at the same time (OperationAborted), and
// Pulumi runs independent resources in parallel, so chain them one after another.
const ownership = new aws.s3.BucketOwnershipControls(
  'site',
  { bucket: bucket.id, rule: { objectOwnership: 'BucketOwnerEnforced' } },
  opts,
);
const publicAccessBlock = new aws.s3.BucketPublicAccessBlock(
  'site',
  {
    bucket: bucket.id,
    blockPublicAcls: true,
    ignorePublicAcls: true,
    blockPublicPolicy: true,
    restrictPublicBuckets: true,
  },
  { provider, dependsOn: [ownership] },
);
const encryption = new aws.s3.BucketServerSideEncryptionConfiguration(
  'site',
  {
    bucket: bucket.id,
    rules: [
      {
        applyServerSideEncryptionByDefault: { sseAlgorithm: 'AES256' },
        bucketKeyEnabled: true,
      },
    ],
  },
  { provider, dependsOn: [publicAccessBlock] },
);
// Versioning gives a way back from a bad upload; old versions expire after 30 days.
const versioning = new aws.s3.BucketVersioning(
  'site',
  { bucket: bucket.id, versioningConfiguration: { status: 'Enabled' } },
  { provider, dependsOn: [encryption] },
);
const lifecycle = new aws.s3.BucketLifecycleConfiguration(
  'site',
  {
    bucket: bucket.id,
    rules: [
      {
        id: 'expire-old-versions',
        status: 'Enabled',
        filter: {},
        noncurrentVersionExpiration: { noncurrentDays: 30 },
        abortIncompleteMultipartUpload: { daysAfterInitiation: 7 },
      },
    ],
  },
  { provider, dependsOn: [versioning] },
);

// --- ACM certificate, validated through the Route 53 zone ---------------------------
const certificate = new aws.acm.Certificate(
  'site',
  { domainName: domain, validationMethod: 'DNS' },
  { provider: usEast1 },
);
const validationOption = certificate.domainValidationOptions.apply(
  (options) => {
    const match = options.find((o) => o.domainName === domain);
    if (!match) throw new Error(`ACM returned no DNS validation for ${domain}`);
    return match;
  },
);
const validationRecord = new aws.route53.Record(
  'cert-validation',
  {
    zoneId: zone.zoneId,
    name: validationOption.resourceRecordName,
    type: validationOption.resourceRecordType,
    records: [validationOption.resourceRecordValue],
    ttl: 300,
    allowOverwrite: true,
  },
  opts,
);
const validation = new aws.acm.CertificateValidation(
  'site',
  {
    certificateArn: certificate.arn,
    validationRecordFqdns: [validationRecord.fqdn],
  },
  { provider: usEast1 },
);

// --- CloudFront ---------------------------------------------------------------------
const originAccessControl = new aws.cloudfront.OriginAccessControl(
  'site',
  {
    name: 'shell-explain-site',
    description: 'CloudFront access to the private site bucket',
    originAccessControlOriginType: 's3',
    signingBehavior: 'always',
    signingProtocol: 'sigv4',
  },
  opts,
);

// Directory index and trailing-slash redirect (see the file's header).
const viewerRequest = new aws.cloudfront.Function(
  'viewer-request',
  {
    name: 'shell-explain-viewer-request',
    runtime: 'cloudfront-js-2.0',
    comment: 'directory index and trailing-slash redirect',
    publish: true,
    code: fs.readFileSync(
      path.join(__dirname, 'functions', 'viewer-request.js'),
      'utf8',
    ),
  },
  opts,
);

// Scripts and styles all come from this host. The web build must not inline scripts
// (web/astro.config.mjs sets assetsInlineLimit to 0), so script-src stays 'self'.
const contentSecurityPolicy = [
  "default-src 'self'",
  "script-src 'self'",
  "style-src 'self'",
  "img-src 'self' data:",
  "font-src 'self'",
  "connect-src 'self'",
  "object-src 'none'",
  "base-uri 'self'",
  "form-action 'self'",
  "frame-ancestors 'none'",
].join('; ');

const headersPolicy = new aws.cloudfront.ResponseHeadersPolicy(
  'security-headers',
  {
    name: 'shell-explain-security-headers',
    comment: `Security headers for ${domain}`,
    securityHeadersConfig: {
      strictTransportSecurity: {
        accessControlMaxAgeSec: 31536000,
        includeSubdomains: false,
        preload: false,
        override: true,
      },
      contentSecurityPolicy: {
        contentSecurityPolicy,
        override: true,
      },
      contentTypeOptions: { override: true },
      frameOptions: { frameOption: 'DENY', override: true },
      referrerPolicy: {
        referrerPolicy: 'strict-origin-when-cross-origin',
        override: true,
      },
    },
    customHeadersConfig: {
      items: [
        {
          header: 'Permissions-Policy',
          value: 'camera=(), microphone=(), geolocation=()',
          override: true,
        },
      ],
    },
  },
  opts,
);

// AWS-managed cache policy: compression on, and the origin's Cache-Control decides the
// TTL (immutable for hashed files, short for HTML; see the deploy workflow).
const cachingOptimized = aws.cloudfront.getCachePolicyOutput(
  { name: 'Managed-CachingOptimized' },
  { provider },
);

const originId = 'site-bucket';
const distribution = new aws.cloudfront.Distribution(
  'site',
  {
    enabled: true,
    comment: domain,
    aliases: [domain],
    httpVersion: 'http2and3',
    isIpv6Enabled: true,
    priceClass: 'PriceClass_All',
    origins: [
      {
        originId,
        domainName: bucket.bucketRegionalDomainName,
        originAccessControlId: originAccessControl.id,
      },
    ],
    defaultCacheBehavior: {
      targetOriginId: originId,
      viewerProtocolPolicy: 'redirect-to-https',
      allowedMethods: ['GET', 'HEAD', 'OPTIONS'],
      cachedMethods: ['GET', 'HEAD'],
      cachePolicyId: cachingOptimized.id,
      responseHeadersPolicyId: headersPolicy.id,
      compress: true,
      functionAssociations: [
        { eventType: 'viewer-request', functionArn: viewerRequest.arn },
      ],
    },
    // With origin access control and no s3:ListBucket, S3 answers 403 for a missing
    // key, so both 403 and 404 serve the site's own 404 page with a real 404 status.
    customErrorResponses: [403, 404].map((errorCode) => ({
      errorCode,
      responseCode: 404,
      responsePagePath: '/404.html',
      errorCachingMinTtl: 10,
    })),
    restrictions: { geoRestriction: { restrictionType: 'none' } },
    viewerCertificate: {
      acmCertificateArn: validation.certificateArn,
      sslSupportMethod: 'sni-only',
      minimumProtocolVersion: 'TLSv1.2_2021',
    },
  },
  opts,
);

// Only this distribution may read the bucket.
new aws.s3.BucketPolicy(
  'site',
  {
    bucket: bucket.id,
    policy: pulumi.jsonStringify({
      Version: '2012-10-17',
      Statement: [
        {
          Sid: 'AllowCloudFrontRead',
          Effect: 'Allow',
          Principal: { Service: 'cloudfront.amazonaws.com' },
          Action: 's3:GetObject',
          Resource: pulumi.interpolate`${bucket.arn}/*`,
          Condition: { StringEquals: { 'AWS:SourceArn': distribution.arn } },
        },
        {
          Sid: 'DenyInsecureTransport',
          Effect: 'Deny',
          Principal: '*',
          Action: 's3:*',
          Resource: [bucket.arn, pulumi.interpolate`${bucket.arn}/*`],
          Condition: { Bool: { 'aws:SecureTransport': 'false' } },
        },
      ],
    }),
  },
  { provider, dependsOn: [lifecycle] },
);

// Alias records (IPv4 and IPv6). Z2FDTNDATAQYW2 is the fixed hosted zone ID that every
// CloudFront distribution uses.
for (const type of ['A', 'AAAA']) {
  new aws.route53.Record(
    `alias-${type}`,
    {
      zoneId: zone.zoneId,
      name: domain,
      type,
      aliases: [
        {
          name: distribution.domainName,
          zoneId: 'Z2FDTNDATAQYW2',
          evaluateTargetHealth: false,
        },
      ],
    },
    opts,
  );
}

// --- Deploy role for GitHub Actions (OIDC, no stored keys) --------------------------
// Trusts only this repository's `production` environment. It can sync the site bucket
// and invalidate this one distribution, nothing else.
const deployRole = new aws.iam.Role(
  'deploy',
  {
    name: 'shell-explain-deploy',
    description: `GitHub Actions deploy of ${domain} (production environment only)`,
    maxSessionDuration: 3600,
    assumeRolePolicy: pulumi.jsonStringify({
      Version: '2012-10-17',
      Statement: [
        {
          Effect: 'Allow',
          Principal: { Federated: oidcProviderArn },
          Action: 'sts:AssumeRoleWithWebIdentity',
          Condition: {
            StringEquals: {
              'token.actions.githubusercontent.com:aud': 'sts.amazonaws.com',
              'token.actions.githubusercontent.com:sub': `${githubSubjectPrefix}:environment:${githubEnvironment}`,
            },
          },
        },
      ],
    }),
  },
  opts,
);
new aws.iam.RolePolicy(
  'deploy',
  {
    role: deployRole.id,
    policy: pulumi.jsonStringify({
      Version: '2012-10-17',
      Statement: [
        {
          Sid: 'ListSiteBucket',
          Effect: 'Allow',
          Action: [
            's3:ListBucket',
            's3:GetBucketLocation',
            's3:ListBucketMultipartUploads',
          ],
          Resource: bucket.arn,
        },
        {
          Sid: 'WriteSiteObjects',
          Effect: 'Allow',
          Action: [
            's3:PutObject',
            's3:GetObject',
            's3:DeleteObject',
            's3:AbortMultipartUpload',
          ],
          Resource: pulumi.interpolate`${bucket.arn}/*`,
        },
        {
          Sid: 'InvalidateDistribution',
          Effect: 'Allow',
          Action: ['cloudfront:CreateInvalidation', 'cloudfront:GetInvalidation'],
          Resource: distribution.arn,
        },
      ],
    }),
  },
  opts,
);

export const distributionId = distribution.id;
export const distributionDomain = distribution.domainName;
export const deployRoleArn = deployRole.arn;
export const url = `https://${domain}/`;
