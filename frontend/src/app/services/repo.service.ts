import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Repo, Response } from '../models/repo.model';

@Injectable({
  providedIn: 'root'
})
export class RepoService {
  private apiUrl = 'http://localhost:4000'; // This will be proxied to your Encore backend

  constructor(private http: HttpClient) { }

  get(id: string): Observable<Repo> {
    return this.http.get<Repo>(`${this.apiUrl}/repos/${id}`);
  }

  list(): Observable<Response> {
    return this.http.get<Response>(`${this.apiUrl}/repos`);
  }

  add(owner: string, name: string): Observable<Repo> {
    return this.http.post<Repo>(`${this.apiUrl}/repos`, { owner, name });
  }

  delete(id: string): Observable<any> {
    return this.http.delete(`${this.apiUrl}/repos/${id}`);
  }
}
