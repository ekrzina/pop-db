"use client"
import { createPerson, getPersons } from "@/api/client"
import BloodTypePickerModal from "@/components/HelperModals"
import { useState } from "react"

type MedicalData = {
    height: number
    weight: number
    bloodType: string
    medicalConditions: string
}

type PersonForm = {
    name: string
    surname: string
    occupation: string
    nationality: string
    city: string
    notes: string
    dateOfBirth: string
    medical: MedicalData
}

export default function CreatePersonModal({
    onClose,
    onSave,
}: {
    onClose: () => void
    onSave: (newPerson: PersonForm & { id: number }) => void
}) {
    const [form, setForm] = useState<PersonForm>({
        name: "",
        surname: "",
        occupation: "",
        nationality: "",
        city: "",
        notes: "",
        dateOfBirth: new Date().toISOString().split("T")[0],
        medical: { height: 0, weight: 0, bloodType: "A+", medicalConditions: "" },
    })
    const [showBloodTypePicker, setShowBloodTypePicker] = useState(false)
    const handleChange = (field: keyof Omit<PersonForm, "medical">, value: string) => {
        setForm((prev) => ({ ...prev, [field]: value }))
    }

    const handleMedicalChange = (field: keyof MedicalData, value: string | number) => {
        setForm((prev) => ({
            ...prev,
            medical: { ...prev.medical, [field]: value },
        }))
    }

    const handleSubmit = async () => {
        try {
            await createPerson(form)
            const updatedList = await getPersons()
            onSave(updatedList[updatedList.length - 1])
            onClose()
        } catch (err) {
            console.error("Failed to create person:", err)
            alert("Failed to create person.")
        }
    }

    const fields: (keyof Omit<PersonForm, "medical">)[] = [
        "name",
        "surname",
        "occupation",
        "dateOfBirth",
        "nationality",
        "city",
        "notes",
    ]

    return (
        <div className="fixed inset-0 bg-black/40 backdrop-blur-sm flex justify-center items-center p-4">
            <div className="bg-white rounded-2xl shadow-xl w-full max-w-xl max-h-[90vh] flex flex-col">

                {/* Header */}
                <div className="flex justify-between items-center px-5 py-3 border-b">
                    <h2 className="text-lg font-semibold">Create Person</h2>
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
                    {fields.map((f) => (
                        <div key={f}>
                            <label className="block text-sm font-medium capitalize mb-1">{f}</label>
                            <input
                                type={f === "dateOfBirth" ? "date" : "text"}
                                className="border p-2 w-full rounded-lg bg-gray-50 focus:ring-2 focus:ring-rose-400 outline-none"
                                value={form[f] || ""}
                                onChange={(e) => handleChange(f, e.target.value)}
                            />
                        </div>
                    ))}

                    {/* Medical info */}
                    <div className="border-t pt-4 space-y-3">
                        <div className="text-sm font-semibold">Medical Data</div>

                        <input
                            type="number"
                            value={form.medical.height}
                            onChange={(e) => handleMedicalChange("height", Number(e.target.value))}
                            placeholder="Height (cm)"
                            className="border p-2 w-full rounded-lg bg-gray-50 focus:ring-2 focus:ring-rose-400 outline-none"
                        />

                        <input
                            type="number"
                            value={form.medical.weight}
                            onChange={(e) => handleMedicalChange("weight", Number(e.target.value))}
                            placeholder="Weight (kg)"
                            className="border p-2 w-full rounded-lg bg-gray-50 focus:ring-2 focus:ring-rose-400 outline-none"
                        />

                        <button
                            type="button"
                            className="border p-2 w-full rounded-lg bg-gray-50 text-left hover:bg-gray-100 transition"
                            onClick={() => setShowBloodTypePicker(true)}
                        >
                            {form.medical.bloodType || "Select Blood Type"}
                        </button>

                        <textarea
                            value={form.medical.medicalConditions}
                            onChange={(e) => handleMedicalChange("medicalConditions", e.target.value)}
                            placeholder="Medical Conditions"
                            className="border p-2 w-full rounded-lg bg-gray-50 focus:ring-2 focus:ring-rose-400 outline-none"
                        />
                    </div>

                    {showBloodTypePicker && (
                        <BloodTypePickerModal
                            value={form.medical.bloodType}
                            onSelect={(bt) => {
                                handleMedicalChange("bloodType", bt)
                                setShowBloodTypePicker(false)
                            }}
                            onClose={() => setShowBloodTypePicker(false)}
                        />
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
                        Create
                    </button>
                </div>
            </div>
        </div>
    )
}