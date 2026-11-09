import { defineCollection } from 'astro:content';
import { z } from 'astro/zod';
import { file } from 'astro/loaders';

const section = z.object({
  id: z.string(),
  implementation: z.string(),
  title: z.string(),
  options: z.array(z.string()),
  markdown: z.string(),
  html: z.string(),
  code: z.array(z.string()),
  example_text: z.string().optional(),
  option_text: z.string().optional(),
});
const commands = defineCollection({
  loader: file('src/data/catalog.json', {
    parser: (text) => JSON.parse(text).entries,
  }),
  schema: z.object({
    schema_version: z.literal(1),
    title: z.string(),
    summary: z.string(),
    aliases: z.array(z.string()).nullable(),
    tags: z.array(z.string()),
    implementations: z.array(z.string()),
    updated: z.string().regex(/^\d{4}-\d{2}-\d{2}$/),
    related_guides: z.array(z.string()).optional(),
    references: z.array(z.url()),
    version_policy: z.string(),
    sections: z.array(section),
  }),
});
export const collections = { commands };
