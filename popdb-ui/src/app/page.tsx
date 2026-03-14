"use client"

import { deletePerson, getPersons, searchPersons } from "@/api/client"
import CreatePersonModal from "@/components/CreatePersonModal"
import EditPersonModal from "@/components/EditPersonModal"
import PersonCard from "@/components/PersonCard"
import SearchBar from "@/components/SearchBar"
import { useEffect, useState } from "react"

export default function Home() {
  const [persons, setPersons] = useState<any[]>([])
  const [selectedIndex, setSelectedIndex] = useState(0)
  const [searchField, setSearchField] = useState<"name" | "surname" | "occupation">("name")
  const [searchQuery, setSearchQuery] = useState("")
  const [page, setPage] = useState(1)
  const [hasNext, setHasNext] = useState(false)
  const [editingPerson, setEditingPerson] = useState<any | null>(null)
  const [creatingPerson, setCreatingPerson] = useState(false)
  const [loading, setLoading] = useState(false)

  const LIMIT = 100
  const offset = (page - 1) * LIMIT

  const fetchData = async () => {
    try {
      const fetchLimit = LIMIT + 1
      const data = searchQuery
        ? await searchPersons(searchField, searchQuery, fetchLimit, offset)
        : await getPersons(fetchLimit, offset)

      const slice = data.length > LIMIT ? data.slice(0, LIMIT) : data
      setPersons(slice)
      setSelectedIndex(0)
      setHasNext(data.length > LIMIT)
    } catch (error) {
      console.error("Failed to fetch persons:", error)
    }
  }

  useEffect(() => {
    if (typeof window === "undefined") return
    fetchData()
  }, [searchField, searchQuery, page])

  // Keyboard navigation
  useEffect(() => {
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === "ArrowRight") nextPerson()
      if (e.key === "ArrowLeft") prevPerson()
    }

    window.addEventListener("keydown", handleKey)
    return () => window.removeEventListener("keydown", handleKey)
  }, [persons, selectedIndex])

  const currentPerson = persons[selectedIndex] || null

  const nextPerson = () => {
    if (!persons.length) return
    setSelectedIndex((prev) => Math.min(prev + 1, persons.length - 1))
  }

  const prevPerson = () => {
    if (!persons.length) return
    setSelectedIndex((prev) => Math.max(prev - 1, 0))
  }

  const nextPage = () => {
    if (hasNext) setPage((prev) => prev + 1)
  }

  const prevPage = () => {
    setPage((prev) => Math.max(1, prev - 1))
  }

  const handleSearch = (field: "name" | "surname" | "occupation", query: string) => {
    setSearchField(field)
    setSearchQuery(query)
    setPage(1)
    setSelectedIndex(0)
  }

  const handleDelete = async () => {
    if (!currentPerson || loading) return
    if (!confirm(`Are you sure you want to delete ${currentPerson.name} ${currentPerson.surname}?`)) return
    try {
      setLoading(true)
      await deletePerson(currentPerson.id)
      await fetchData()
    } finally {
      setLoading(false)
    }
  }

  const handleEdit = () => {
    if (!currentPerson) return
    setEditingPerson(currentPerson)
  }

  const handleSaveEdit = async () => {
    setEditingPerson(null)
    setCreatingPerson(false)
    await fetchData()
  }

  const handleCreate = () => setCreatingPerson(true)

  return (
    // Screen p4, larger screens p8 or p10, max width 4xl, centered
    <main className="p-4 sm:p-8 lg:p-10 space-y-6 max-w-4xl mx-auto">
      <SearchBar onSearch={handleSearch} />

      {/* Scrollable list of current page */}
      {persons.length > 0 && (
        <div className="flex gap-2 overflow-x-auto pb-2 scrollbar-thin">
          {persons.map((p, idx) => (
            <button
              key={p.id}
              onClick={() => setSelectedIndex(idx)}
              className={`border px-3 py-1 rounded whitespace-nowrap ${persons[selectedIndex]?.id === p.id ? "bg-black text-white" : ""}`}
            >
              {p.name} {p.surname}
            </button>
          ))}
        </div>
      )}

      <div className="flex gap-2 items-center text-sm text-gray-700">
        <span>Page {page}</span>
        <button
          className="px-3 py-1 rounded bg-white shadow hover:shadow-md"
          onClick={prevPage}
          disabled={page === 1}
        >
          Prev page
        </button>
        <button
          className="px-3 py-1 rounded bg-white shadow hover:shadow-md"
          onClick={nextPage}
          disabled={!hasNext}
        >
          Next page
        </button>
      </div>

      {/* Navigation buttons */}
      <div className="flex flex-wrap gap-3 items-center">
        <button
          className="px-5 py-2.5 rounded-xl bg-white shadow hover:shadow-md"
          onClick={prevPerson}
          disabled={selectedIndex === 0}
        >
          Previous
        </button>
        <button
          className="px-5 py-2.5 rounded-xl bg-white shadow hover:shadow-md"
          onClick={nextPerson}
          disabled={selectedIndex >= persons.length - 1}
        >
          Next
        </button>

        <button
          className="px-5 py-2.5 rounded-xl bg-white shadow hover:shadow-md"
          onClick={handleEdit}
          disabled={!currentPerson}
        >
          ✏️ Edit
        </button>
        <button
          className="px-5 py-2.5 rounded-xl bg-white shadow hover:shadow-md"
          onClick={handleDelete}
          disabled={!currentPerson || loading}
        >
          {loading ? "Deleting..." : "🗑️ Delete"}
        </button>

        <button
          onClick={handleCreate}
          className="px-4 py-2 rounded-xl bg-rose-500 text-white shadow 
                transition-all duration-300 
                hover:bg-orange-500"
        >
          ✚ Add Person
        </button>

        <a
          href="/summary"
          target="_blank"
          rel="noopener noreferrer"
          className="border px-4 py-2 rounded inline-block hover:bg-gray-100 transition"
        >
          📋 Summary ↗
        </a>
      </div>

      {/* Current person card */}
      <div className="flex justify-center">
        {currentPerson ? <PersonCard person={currentPerson} /> : <div className="text-center text-gray-500">
          No persons found
        </div>}
      </div>

      {editingPerson && (
        <EditPersonModal
          person={editingPerson}
          onClose={() => setEditingPerson(null)}
          onSave={handleSaveEdit}
        />
      )}
      {creatingPerson && (
        <CreatePersonModal
          onClose={() => setCreatingPerson(false)}
          onSave={async () => {
            setCreatingPerson(false)
            await fetchData()
          }}
        />
      )}
    </main>
  )
}