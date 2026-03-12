"use client"

import { useState, useEffect, useRef } from "react"

type SearchField = "name" | "surname" | "occupation"

export default function SearchBar({
    onSearch,
}: {
    onSearch: (field: SearchField, query: string) => void
}) {
    const [field, setField] = useState<SearchField>("name")
    const [query, setQuery] = useState("")
    const inputRef = useRef<HTMLInputElement>(null)

    const submit = () => {
        onSearch(field, query)
    }

    // Keyboard shortcut: press "/" to focus search
    useEffect(() => {
        const handler = (e: KeyboardEvent) => {
            if (e.key === "/") {
                e.preventDefault()
                inputRef.current?.focus()
            }
        }

        window.addEventListener("keydown", handler)
        return () => window.removeEventListener("keydown", handler)
    }, [])

    const clearSearch = () => {
        setQuery("")
        onSearch(field, "")
        inputRef.current?.focus()
    }

    return (
        <div className="flex items-center gap-3 bg-white shadow-md rounded-2xl px-4 py-3 focus-within:ring-2 focus-within:ring-indigo-400">

            {/* Field selector */}
            <select
                value={field}
                onChange={(e) => setField(e.target.value as SearchField)}
                className="bg-gray-50 px-2 py-1 rounded-lg text-gray-600 outline-none"
            >
                <option value="name">Name</option>
                <option value="surname">Surname</option>
                <option value="occupation">Occupation</option>
            </select>

            {/* Search input */}
            <div className="flex items-center gap-2 flex-1 text-gray-400">
                <span>🔍</span>

                <input
                    ref={inputRef}
                    className="flex-1 bg-transparent outline-none text-gray-700"
                    placeholder="Search person..."
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    onKeyDown={(e) => e.key === "Enter" && submit()}
                />
            </div>

            {/* Clear button */}
            {query && (
                <button
                    onClick={clearSearch}
                    className="text-gray-400 hover:text-gray-600"
                >
                    ✕
                </button>
            )}

            {/* Search button */}
            <button 
                onClick={submit}
                className="px-4 py-2 rounded-xl bg-rose-500 text-white shadow
             transition-colors duration-300
             hover:bg-orange-500"
            >
                Search
            </button>
        </div>
    )
}