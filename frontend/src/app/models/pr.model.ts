export interface PullRequest {
  id: string;
  number: number;
  title: string;
  author: string;
  avatar_url: string;
  html_url: string;
  labels: Label[];
  status: string;
  draft: boolean;
  created_at: string;
}

export interface Label {
  name: string;
  color: string;
}

export interface Response {
  pull_requests: PullRequest[];
}
