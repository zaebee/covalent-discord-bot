import axios from 'axios'


const baseURL = process.env.ELASTICSEARCH_URL

const axiosConfig = {
  baseURL: baseURL,
  timeout: 30000,
  headers: {'X-Requested-With': 'XMLHttpRequest'},
}

let $http = axios.create(axiosConfig)

export default $http
