import { AutocompleteOptions } from "@algolia/autocomplete-js";
import { SearchOptions } from "@algolia/client-search";
import { SearchClient } from "algoliasearch/lite";
import { FC } from "react";
type AlgoliaRecord = {
    objectID: string;
    url: string;
    origin: string;
    title: string;
    content: string;
    lang?: string;
    description?: string;
    keywords?: string[];
    image?: string;
    authors?: string[];
    datePublished?: number;
    dateModified?: number;
    category?: string;
    hierarchy?: Hierarchy;
    hierarchicalCategories?: Hierarchy;
    urlDepth?: number;
    position?: number;
};
type Hierarchy = {
    [lvl: string]: string;
};
interface AlgoliaAutocompleteProps extends Partial<AutocompleteOptions<AlgoliaRecord>> {
    branch: string;
    debounceTimeout?: number;
    searchClient: SearchClient;
    searchOptions?: SearchOptions;
    siteId: string;
}
declare const AlgoliaAutocomplete: FC<AlgoliaAutocompleteProps>;
export type { AlgoliaAutocompleteProps };
export { AlgoliaAutocomplete };
