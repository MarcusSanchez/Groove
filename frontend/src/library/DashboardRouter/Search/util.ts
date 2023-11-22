import { Album, Artist, DataType, SearchType, Track } from "./types.ts";

export function redirect(id: string, type: SearchType, navigate: any) {
  switch (type) {
    case SearchType.Album:
      navigate(`/dashboard/${type}?id=${id}`);
      break;
    case SearchType.Artist:
      navigate(`/dashboard/artist?id=${id}`);
      break;
    case SearchType.Track:
      navigate(`/dashboard/track?id=${id}`);
      break;
  }
}


export function sortResults(Artists?: Artist[], Albums?: Album[], Tracks?: Track[]): DataType[] {
  const artists = Artists?.map((artist) => ({ data: artist, type: SearchType.Artist })) || [] as DataType[];
  const albums = Albums?.map((album) => ({ data: album, type: SearchType.Album })) || [] as DataType[];
  const tracks = Tracks?.map((track) => ({ data: track, type: SearchType.Track })) || [] as DataType[];

  const randomChoice = (arr: (DataType | undefined)[]) => {
    let randomIndex = Math.floor(Math.random() * arr.length);
    return arr.splice(randomIndex, 1)[0];
  }

  let sortedResults = [] as DataType[];
  for (let i = 0; i < 3; i++) {
    let arr = [artists.shift(), albums.shift(), tracks.shift()];

    while (arr.length > 0) {
      let choice = randomChoice(arr);
      if (choice) sortedResults.push(choice);
    }
  }

  while (tracks.length > 0 || albums.length > 0) {
    let arr = [tracks.shift(), albums.shift()];

    while (arr.length > 0) {
      let choice = randomChoice(arr);
      if (choice) sortedResults.push(choice);
    }
  }

  sortedResults = [...sortedResults, ...artists];
  return sortedResults;
}

export function msToMinutesSeconds(durationMs: number): string {
  const durationSec: number = Math.floor(durationMs / 1000);
  const minutes: number = Math.floor(durationSec / 60);
  const seconds: number = durationSec % 60;

  // Ensure seconds are displayed with leading zero if needed
  return `${minutes}:${seconds.toString().padStart(2, '0')}`;
}

