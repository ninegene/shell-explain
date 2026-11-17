# Hosting infrastructure

shell-explain.ninegene.com is the static Astro site (`web/dist/`) in a private S3 bucket behind CloudFront. The Pulumi program in `pulumi/` creates everything specific to this site. It follows [ninegene/ninegene.com](https://github.com/ninegene/ninegene.com) (`infra/`), which owns the shared parts.

| Piece | Where it comes from |
| --- | --- |
| Pulumi state bucket, KMS secrets key, GitHub OIDC provider | Created once by ninegene.com's `infra/bootstrap/bootstrap.sh`. Shared, not managed here. |
| Route 53 zone `ninegene.com` | Owned and protected by the ninegene.com stack. This program only looks it up and adds records. |
| ACM certificate (`us-east-1`), validation record | This program |
| Site bucket `shell-explain-site`, CloudFront distribution, function, security headers | This program |
| Alias records `A` and `AAAA` for `shell-explain.ninegene.com` | This program |
| IAM role `shell-explain-deploy` (GitHub OIDC, `production` environment only) | This program |

The AWS account is a project account from AWS's new sign-up with advanced features active. Regional resources go in `us-east-2`; only global services and the CloudFront certificate use `us-east-1`. Sign in with `aws login --profile ninegene`.

Never commit the AWS account ID, keys or `.env` files. The account ID is a GitHub environment secret.

Cost: a few cents a month (the hosted zone is shared and already paid for).

## Check before any change

```bash
export AWS_PROFILE=ninegene
aws sts get-caller-identity --query Account --output text
```

The account ID must match the `ninegene` project account listed in your global agent rules. Stop if it does not.

## First setup

Each step that writes to AWS needs a deliberate go-ahead.

```bash
just infra-setup
cd infra/pulumi
pulumi stack init prod --secrets-provider "awskms://alias/ninegene-pulumi?region=us-east-2"
pulumi config set region us-east-2
pulumi config set githubSubjectPrefix "$(gh api repos/ninegene/shell-explain/actions/oidc/customization/sub --jq .sub_claim_prefix)"
pulumi preview
```

The repository uses immutable OIDC subjects, so the deploy role's trust policy needs the `sub_claim_prefix` from that API call (it contains numeric repository IDs, not secrets). Check the preview: it must only create resources, and must change nothing in the ninegene.com stack. Then `pulumi up`. Certificate validation takes a few minutes.

Commit the generated `Pulumi.prod.yaml` (it holds the encrypted stack key, no secrets).

## GitHub environment

In Settings > Environments, create `production` and set these variables from the stack outputs (`pulumi stack output`):

- `SITE_URL`: `https://shell-explain.ninegene.com`
- `SITE_BUCKET`: `siteBucket`
- `CLOUDFRONT_DISTRIBUTION_ID`: `distributionId`
- `AWS_REGION`: optional, defaults to `us-east-2`

and the secret `AWS_ACCOUNT_ID`.

## Deploying

`.github/workflows/deploy.yml` builds `web/` and publishes it after "Native Linux AMD64" succeeds on a push to `main`. To deploy or roll back by hand, run the workflow from the Actions tab with a branch, tag or full commit SHA.

## Local checks

```bash
just infra-check     # typecheck and CloudFront function tests
just infra-preview   # pulumi preview (read-only)
```

## Content-Security-Policy

The CloudFront response headers allow scripts and styles only from the site itself. `web/astro.config.mjs` sets `assetsInlineLimit: 0` so Astro never inlines a script into a page. After a web change, check that `grep -rhoE '<script[^>]*>' web/dist --include='*.html'` shows only `src=` scripts.
