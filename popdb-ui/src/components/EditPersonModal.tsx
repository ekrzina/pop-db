"use client"
import { createPerson, updatePerson } from "@/api/client"
import { useQueryClient } from "@tanstack/react-query"
import { useState } from "react"

export default function EditPersonModal({
    person,
    onClose,
    onSave,
    isNew = false, // flag to differentiate create vs edit
}: {
    person: any
    onClose: () => void
    onSave: (updated: any) => void
    isNew?: boolean
}) {
    const [form, setForm] = useState({ ...person })
    const queryClient = useQueryClient()
    const handleChange = (field: string, value: string) => {
        setForm((prev: typeof form) => ({ ...prev, [field]: value }))
    }

    const handleMedicalChange = (field: string, value: string | number) => {
        setForm((prev: typeof form) => ({
            ...prev,
            medical: { ...prev.medical, [field]: value },
        }))
    }

    const handleSubmit = async () => {
        try {
            const updatedData = { ...form }
            if (updatedData.medical) {
                updatedData.medical = {
                    height: updatedData.medical.height,
                    weight: updatedData.medical.weight,
                    bloodType: updatedData.medical.bloodType,
                    medicalConditions: updatedData.medical.medicalConditions,
                }
            }

            let saved
            if (isNew) {
                saved = await createPerson(updatedData)
            } else {
                saved = await updatePerson(person.id, updatedData)
                saved = { ...person, ...updatedData }
            }
            queryClient.invalidateQueries({ queryKey: ["person", person.id] })
            onSave(saved)
            onClose()
        } catch (err) {
            console.error("Failed to save person:", err)
            alert("Failed to save person.")
        }
    }

    return (
        <div className="fixed inset-0 bg-black bg-opacity-30 flex justify-center items-center">
            <div className="bg-white p-6 rounded shadow-lg w-96 space-y-4">
                <h2 className="text-xl font-semibold">{isNew ? "Create Person" : "Edit Person"}</h2>

                {/* Basic info */}
                {["name", "surname", "occupation", "dateOfBirth", "nationality", "city", "notes"].map((f) => (
                    <div key={f}>
                        <label className="block text-sm font-medium capitalize">{f}</label>
                        <input
                            type="text"
                            className="border p-2 w-full rounded"
                            value={form[f] || ""}
                            onChange={(e) => handleChange(f, e.target.value)}
                        />
                    </div>
                ))}

                {/* Medical info */}
                {form.medical && (
                    <div className="border-t pt-4 space-y-2">
                        <div className="text-sm font-medium">Medical Data</div>

                        <div>
                            <label className="block text-xs font-medium">Height (cm)</label>
                            <input
                                type="number"
                                className="border p-2 w-full rounded"
                                value={form.medical.height || ""}
                                onChange={(e) => handleMedicalChange("height", Number(e.target.value))}
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-medium">Weight (kg)</label>
                            <input
                                type="number"
                                className="border p-2 w-full rounded"
                                value={form.medical.weight || ""}
                                onChange={(e) => handleMedicalChange("weight", Number(e.target.value))}
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-medium">Blood Type</label>
                            <input
                                type="text"
                                className="border p-2 w-full rounded bg-gray-100"
                                value={form.medical.bloodType || ""}
                                readOnly
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-medium">Medical Conditions</label>
                            <textarea
                                className="border p-2 w-full rounded"
                                value={form.medical.medicalConditions || ""}
                                onChange={(e) => handleMedicalChange("medicalConditions", e.target.value)}
                            />
                        </div>
                    </div>
                )}

                <div className="flex gap-2 justify-end">
                    <button className="px-4 py-2 border rounded" onClick={onClose}>
                        Cancel
                    </button>
                    <button className="px-4 py-2 bg-black text-white rounded" onClick={handleSubmit}>
                        {isNew ? "Create" : "Save"}
                    </button>
                </div>
            </div>
        </div>
    )
}