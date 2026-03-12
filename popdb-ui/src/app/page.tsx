"use client"

import { deletePerson, getPersons } from "@/api/client"
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
  const [editingPerson, setEditingPerson] = useState<any | null>(null)
  const [creatingPerson, setCreatingPerson] = useState(false)

  // Fetch all persons once on client-side
  useEffect(() => {
    if (typeof window === "undefined") return
    getPersons().then((data: any[]) => {
      setPersons(data)
      setSelectedIndex(0) // show first person by default
    })
  }, [])
  // Keyboard navigation
  useEffect(() => {
  const handleKey = (e: KeyboardEvent) => {
    if (e.key === "ArrowRight") {
      nextPerson()
    }

    if (e.key === "ArrowLeft") {
      prevPerson()
    }
  }
  window.addEventListener("keydown", handleKey)
  return () => window.removeEventListener("keydown", handleKey)
}, [persons, selectedIndex, searchQuery])

  // Filter persons only if search query exists
  const filteredPersons = searchQuery
    ? persons.filter((p) => {
      const fieldVal = p[searchField] || ""
      return fieldVal.toLowerCase().includes(searchQuery.toLowerCase())
    })
    : []

  // Current person to show
  const currentPerson = searchQuery
    ? filteredPersons[selectedIndex] || null
    : persons[selectedIndex] || persons[0] || null

  // Navigation
  const nextPerson = () => {
    const list = searchQuery ? filteredPersons : persons
    if (!list.length) return
    setSelectedIndex((prev) => Math.min(prev + 1, list.length - 1))
  }
  const prevPerson = () => {
    const list = searchQuery ? filteredPersons : persons
    if (!list.length) return
    setSelectedIndex((prev) => Math.max(prev - 1, 0))
  }
  // Handle search
  const handleSearch = (field: "name" | "surname" | "occupation", query: string) => {
    setSearchField(field)
    setSearchQuery(query)
    setSelectedIndex(0) // first filtered person
  }
  // Edit / Delete
  const handleDelete = async () => {
    if (!currentPerson) return
    if (!confirm(`Are you sure you want to delete ${currentPerson.name} ${currentPerson.surname}?`)) return
    await deletePerson(currentPerson.id)
    setPersons((prev) => prev.filter((p) => p.id !== currentPerson.id))
    setSelectedIndex((prev) => Math.min(prev, persons.length - 2))
  }
  const handleEdit = () => {
    if (!currentPerson) return
    setEditingPerson(currentPerson)
  }
  const handleSaveEdit = (updated: any, isNew = false) => {
    setPersons((prev) => {
      const updatedList = isNew ? [...prev, updated] : prev.map((p) => (p.id === updated.id ? updated : p))
      if (isNew) setSelectedIndex(updatedList.length - 1)
      return updatedList
    })
    setEditingPerson(null)
    setCreatingPerson(false)
  }
  const handleCreate = () => setCreatingPerson(true)

  return (
    // Screen p4, larger screens p8 or p10, max width 4xl, centered
    <main className="p-4 sm:p-8 lg:p-10 space-y-6 max-w-4xl mx-auto">
      <SearchBar onSearch={handleSearch} />

      {/* Scrollable list only if search is applied */}
      {searchQuery && filteredPersons.length > 0 && (
        <div className="flex gap-2 overflow-x-auto pb-2 scrollbar-thin">
          {filteredPersons.map((p, idx) => (
            <button
              key={p.id}
              onClick={() => {
                setSelectedIndex(persons.findIndex((pp) => pp.id === p.id))
                setSearchQuery("") // reset filter after click
              }}
              className={`border px-3 py-1 rounded whitespace-nowrap ${persons[selectedIndex]?.id === p.id ? "bg-black text-white" : ""
                }`}
            >
              {p.name} {p.surname}
            </button>
          ))}
        </div>
      )}

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
          disabled={!currentPerson}
        >
          🗑️ Delete
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
          onSave={(newPerson) => {
            setPersons((prev) => [...prev, newPerson])
            setSelectedIndex(persons.length)
            setCreatingPerson(false)
          }}
        />
      )}
    </main>
  )
}