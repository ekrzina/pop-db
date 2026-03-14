import axios from "axios"

const api = axios.create({
  baseURL: "http://localhost:8080",
})

export async function getPerson(id: number) {
  const res = await api.get(`/api/v1/persons/${id}`)
  return res.data
}

export async function getPersons(limit = 100, offset = 0) {
  const res = await api.get("/api/v1/persons", {
    params: { limit, offset },
  })
  return res.data
}

export async function deletePerson(id: number) {
  await api.delete(`/api/v1/persons/${id}`)
}

export async function updatePerson(id: number, data: any) {
  await api.put(`/api/v1/persons/${id}`, data)
}

export async function searchPersons(
  field: "name" | "surname" | "occupation",
  query: string,
  limit = 100,
  offset = 0
) {
  const res = await api.get("/api/v1/persons", {
    params: { searchField: field, searchQuery: query, limit, offset },
  })
  return res.data
}

export async function getSummary() {
  const res = await api.get(`/api/v1/persons/summary`)
  return res.data
}

export async function createPerson(data: any) {
  const res = await api.post(`/api/v1/persons`, data)
  return res.data
}