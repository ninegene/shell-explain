import { getCollection } from 'astro:content';
import catalog from '../data/catalog.json';

export const revision = catalog.revision;
export const labels: Record<string, string> = {
  common: 'Portable / common',
  gnu: 'GNU',
  apple: 'macOS / Apple',
  busybox: 'BusyBox',
  bash: 'Bash builtin',
};
export async function commands() {
  return (await getCollection('commands')).sort((a, b) =>
    a.id.localeCompare(b.id),
  );
}
