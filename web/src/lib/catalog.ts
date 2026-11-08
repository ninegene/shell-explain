import { getCollection } from 'astro:content';
import catalog from '../data/catalog.json';

export const revision = catalog.revision;
// Implementation labels come from the Go registry exported with the catalog.
export const labels: Record<string, string> = Object.fromEntries(
  catalog.implementations.map((i: { id: string; label: string }) => [
    i.id,
    i.label,
  ]),
);
export async function commands() {
  return (await getCollection('commands')).sort((a, b) =>
    a.id.localeCompare(b.id),
  );
}
