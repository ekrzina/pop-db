"use client"
import { useState } from "react"

type SearchField = "name" | "surname" | "occupation"

export default function SearchBar({
    onSearch,
}: {
    onSearch: (field: SearchField, query: string) => void
}) {
    const [field, setField] = useState<SearchField>("name")
    const [query, setQuery] = useState("")

    const submit = () => {
        onSearch(field, query)
    }

    return (
        <div className="flex gap-2 mb-4">
            <select
                value={field}
                onChange={(e) => setField(e.target.value as SearchField)}
                className="border p-2"
            >
                <option value="name">Name</option>
                <option value="surname">Surname</option>
                <option value="occupation">Occupation</option>
            </select>

            <input
                className="border p-2 flex-1"
                placeholder="Search..."
                value={query}
                onChange={(e) => setQuery(e.target.value)}
            />

            <button onClick={submit} className="bg-black text-white px-4">
                Search
            </button>
        </div>
    )
}