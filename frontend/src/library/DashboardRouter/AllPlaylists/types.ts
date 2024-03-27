export type Playlists = {
  href: string;
  items: Playlist[];
  limit: number;
  next: null | string;
  offset: number;
  previous: null | string;
  total: number;
};

export type Playlist = {
  collaborative: boolean;
  description: string;
  external_urls: {
    spotify: string;
  };
  href: string;
  id: string;
  images?: Image[];
  name: string;
  owner: {
    display_name: string;
    external_urls: {
      spotify: string;
    };
    href: string;
    id: string;
    type: string;
    uri: string;
  };
  primary_color: null;
  public: boolean;
  snapshot_id: string;
  tracks: {
    href: string;
    total: number;
  };
  type: string;
  uri: string;
};

export type Image = {
  height: number;
  url: string;
  width: number;
};
