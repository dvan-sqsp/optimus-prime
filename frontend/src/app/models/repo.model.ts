export interface Response {
  repos: Repo[];
}

export interface Repo {
  id: string;
  owner: string;
  name: string;
  createdAt?: string;
}
