import axios from "axios";

export const BASE_URL = "http://localhost:8080"
export const AXIOS = axios.create({
  baseURL: "http://localhost:8080"
});
export const VERSIONS = {
  "python" : "3.10.0",
  "go" : "1.16.2"
}
export const SNIPPETS = {
  "python" : "# write your code here",
  "go" : "package main \n\nimport \"fmt\"\n\nfunc main() {\n\t// write your code here \n}"
}