"use client"

import { getSummary } from "@/api/client"
import { useEffect, useState } from "react"

export default function Summary() {
    const [text, setText] = useState("")
    useEffect(() => {
        async function load() {
            const data = await getSummary()
            const formatted = data
                .map((p: any) => `${p.id} ${p.name} ${p.surname}`)
                .join("\n")
            setText(formatted)
        }
        load()
    }, [])
    return (
        <pre style={{ padding: "20px" }}>
            {text}
        </pre>
    )
}