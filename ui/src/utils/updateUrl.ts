import { Params } from '../common/params';
import { SortByFields } from '../store/resource';
import { titleCase } from '../common/titlecase';

// This function returns all selected filters in a combination of params
export const UpdateURL = (
  search: string,
  sort: string,
  categories: string,
  kinds: string,
  catalogs: string
) => {
  // To get URLSearchParams object instance
  const searchParams = new URLSearchParams(window.location.search);

  // To delete query params if search is empty
  if (!search && searchParams.has(Params.Query)) searchParams.delete(Params.Query);

  // Sets query params
  if (search) searchParams.set(Params.Query, search);

  //To delete sort params if doesn't selected any sort
  if (sort === SortByFields.Unknown && searchParams.has(Params.SortBy))
    searchParams.delete(Params.SortBy);

  // Set sort params
  if (sort !== SortByFields.Unknown) searchParams.set(Params.SortBy, sort);

  // To delete category params if doesn't selected any category
  if (!categories && searchParams.has(Params.Category)) searchParams.delete(Params.Category);

  // Sets category params
  if (categories) searchParams.set(Params.Category, categories);

  // To delete kind params if doesn't selected any kind
  if (!kinds && searchParams.has(Params.Kind)) searchParams.delete(Params.Kind);

  // Sets Kind params
  if (kinds) searchParams.set(Params.Kind, kinds);

  // To delete catalog params if doesn't selected any catalog
  if (!catalogs && searchParams.has(Params.Catalog)) searchParams.delete(Params.Catalog);

  // Sets catalos params
  if (catalogs) searchParams.set(Params.Catalog, titleCase(catalogs));

  return searchParams.toString();
};
