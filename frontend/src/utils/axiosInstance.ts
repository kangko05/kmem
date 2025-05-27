import axios from "axios";
import { SERVER } from "../constants";

export const axiosInstance = axios.create({ baseURL: SERVER, withCredentials: true });
