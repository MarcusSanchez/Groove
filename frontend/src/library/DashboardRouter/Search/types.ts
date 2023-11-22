export type Image = {
  height: number;
  url: string;
  width: number;
};

export type ListedArtists = {
  href: string;
  id: string;
  name: string;
  type: string;
};

export type Album = {
  album_type: string;
  artists: ListedArtists[];
  id: string;
  images: Image[];
  name: string;
  release_date: string;
  type: string;
};

export type Track = {
  album: Album;
  artists: ListedArtists[];
  duration_ms: number;
  popularity: number;
  id: string;
  name: string;
  // uri: string;
};

export type Artist = {
  genres: string[];
  id: string;
  images: Image[];
  name: string;
  popularity: number;
  type: string;
}

export type SearchResult = {
  albums?: {
    items: Album[];
  };
  artists?: {
    items: Artist[];
  };
  tracks?: {
    items: Track[];
  };
};

export enum SearchType {
  Album = "album",
  Artist = "artist",
  Track = "track",
}

export interface DataType {
  data: Artist | Album | Track;
  type: SearchType;
}
