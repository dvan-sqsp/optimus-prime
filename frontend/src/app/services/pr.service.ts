import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Response } from '../models/pr.model';

@Injectable({
  providedIn: 'root'
})
export class PRService {
  private apiUrl = 'http://localhost:4000';

  constructor(private http: HttpClient) { }

  list(owner: string, name: string): Observable<Response> {
    return this.http.get<Response>(`${this.apiUrl}/pull_requests/${owner}/${name}`);
  }
}
