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
        // === Validation ===
        if (!form.name.trim() || !form.surname.trim() || !form.city.trim() || !form.nationality.trim()) {
            alert("Name, Surname, City, and Nationality cannot be empty.")
            return
        }

        if (!form.medical || form.medical.height <= 0 || form.medical.weight <= 0) {
            alert("Height and Weight must be greater than 0.")
            return
        }

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
        <div className="fixed inset-0 bg-black/40 backdrop-blur-sm flex justify-center items-center p-4">
            <div className="bg-white rounded-2xl shadow-xl w-full max-w-xl max-h-[90vh] flex flex-col">

                {/* Header */}
                <div className="flex justify-between items-center px-5 py-3 border-b">
                    <h2 className="text-lg font-semibold">
                        {isNew ? "Create Person" : "Edit Person"}
                    </h2>
                    <button
                        onClick={onClose}
                        className="text-gray-500 hover:text-gray-700 text-xl leading-none"
                    >
                        ×
                    </button>
                </div>

                {/* Scrollable content */}
                <div className="p-5 space-y-4 overflow-y-auto">

                    {/* Basic info */}
                    {["name", "surname", "occupation", "dateOfBirth", "nationality", "city", "notes"].map((f) => (
                        <div key={f}>
                            <label className="block text-sm font-medium capitalize mb-1">{f}</label>
                            <input
                                type="text"
                                className="border p-2 w-full rounded-lg bg-gray-50 focus:ring-2 focus:ring-rose-400 outline-none"
                                value={form[f] || ""}
                                onChange={(e) => handleChange(f, e.target.value)}
                            />
                        </div>
                    ))}

                    {/* Medical info */}
                    {form.medical && (
                        <div className="border-t pt-4 space-y-3">
                            <div className="text-sm font-semibold mb-2">Medical Data</div>

                            <div>
                                <label className="block text-sm font-medium mb-1">Height (cm)</label>
                                <input
                                    type="number"
                                    className="border p-2 w-full rounded-lg bg-gray-50 focus:ring-2 focus:ring-rose-400 outline-none"
                                    value={form.medical.height || ""}
                                    onChange={(e) => handleMedicalChange("height", Number(e.target.value))}
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium mb-1">Weight (kg)</label>
                                <input
                                    type="number"
                                    className="border p-2 w-full rounded-lg bg-gray-50 focus:ring-2 focus:ring-rose-400 outline-none"
                                    value={form.medical.weight || ""}
                                    onChange={(e) => handleMedicalChange("weight", Number(e.target.value))}
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium mb-1">Blood Type</label>
                                <input
                                    type="text"
                                    className="border p-2 w-full rounded-lg bg-gray-100"
                                    value={form.medical.bloodType || ""}
                                    readOnly
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium mb-1">Medical Conditions</label>
                                <textarea
                                    className="border p-2 w-full rounded-lg bg-gray-50 focus:ring-2 focus:ring-rose-400 outline-none"
                                    value={form.medical.medicalConditions || ""}
                                    onChange={(e) => handleMedicalChange("medicalConditions", e.target.value)}
                                />
                            </div>
                        </div>
                    )}
                </div>

                {/* Footer */}
                <div className="flex justify-end gap-2 px-5 py-3 border-t">
                    <button
                        className="px-4 py-2 border rounded-lg hover:bg-gray-100"
                        onClick={onClose}
                    >
                        Cancel
                    </button>
                    <button
                        className="px-4 py-2 bg-rose-500 text-white rounded-lg hover:bg-orange-500 transition"
                        onClick={handleSubmit}
                    >
                        {isNew ? "Create" : "Save"}
                    </button>
                </div>
            </div>
        </div>
    )
}