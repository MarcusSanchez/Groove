export type Image = {
  height: number;
  url: string;
  width: number;
};

export type Artist = {
  external_urls: {
    spotify: string
  },
  followers: {
    href: null,
    total: number
  },
  genres: string[],
  href: string,
  id: string,
  images: Image[],
  name: string,
  popularity: number,
  type: string,
  uri: string
}

export type FeaturedArtists = {
  external_urls: {
    spotify: string
  },
  href: string,
  id: string,
  name: string,
  type: string,
  uri: string
}

export type Album = {
  album_type: string,
  artists: FeaturedArtists[],
  external_urls: {
    spotify: string
  },
  href: string,
  id: string,
  images: Image[],
  is_playable: boolean,
  name: string,
  release_date: string,
  release_date_precision: string,
  total_tracks: number,
  type: string,
  uri: string
}

export type Track = {
  album: Album,
  artists: FeaturedArtists[],
  disc_number: number,
  duration_ms: number,
  explicit: boolean,
  external_ids: {
    isrc: string
  },
  external_urls: {
    spotify: string
  },
  href: string,
  id: string,
  is_local: boolean,
  is_playable: boolean,
  name: string,
  popularity: number,
  preview_url: string,
  track_number: number,
  type: string,
  uri: string
}

export type RelatedArtists = {
  artists: Artist[]
}

export type TopTracks = {
  tracks: Track[]
}

export type Albums = {
  href: string,
  items: Album[],
  total: number
}

