import catalog from '../data/catalog.json';

// Guides, series and the glossary come from the same exported catalog as the
// command pages, so the website and `explain learn` show identical content.
export type Guide = (typeof catalog.guides)[number];
export const series = catalog.series;
// Already in reading order: by series, then position.
export const guides: Guide[] = catalog.guides;
export const glossary = catalog.glossary;

export const guideURL = (id: string) => `/learn/${id}/`;
export const glossaryURL = (id?: string) =>
  id ? `/learn/glossary/#${id}` : '/learn/glossary/';

export function findGuide(id: string) {
  return guides.find((g) => g.id === id);
}
export function guideNumber(id: string) {
  return guides.findIndex((g) => g.id === id) + 1;
}
export function guidesIn(seriesId: string) {
  return guides.filter((g) => g.series === seriesId);
}
// Series that have at least one guide; a planned series with none is not shown.
export function activeSeries() {
  return series.filter((s) => guidesIn(s.id).length > 0);
}
export function neighbors(id: string) {
  const i = guides.findIndex((g) => g.id === id);
  return { prev: guides[i - 1], next: guides[i + 1] };
}
export function seriesOf(guide: Guide) {
  return series.find((s) => s.id === guide.series)!;
}
