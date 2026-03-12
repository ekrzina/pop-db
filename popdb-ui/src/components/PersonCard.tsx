"use client"

type MedicalData = {
  height: number
  weight: number
  bloodType: string
  medicalConditions: string
}

type Person = {
  id: number
  name: string
  surname: string
  occupation: string
  nationality: string
  city: string
  notes: string
  dateOfBirth: string
  medical: MedicalData
}

export default function PersonCard({ person }: { person: Person }) {
  return (
    <div className="max-w-2xl w-full bg-white rounded-3xl shadow-lg p-8 space-y-4">
      <h2 className="text-2xl font-semibold text-[var(--foreground)]">
        {person.name} {person.surname}
      </h2>
      <div>City: {person.city}</div>
      <div>Nationality: {person.nationality}</div>
      <div>Date of Birth: {person.dateOfBirth}</div>
      <div>Occupation: {person.occupation}</div>

      <div>
        <div className="font-medium">Notes</div>
        <div className="border rounded p-2 max-h-32 overflow-y-auto bg-[var(--card-bg)] whitespace-pre-line">
          {person.notes || "None"}
        </div>
      </div>

      {person.medical && (
        <>
          <div>Height: {person.medical.height} cm</div>
          <div>Weight: {person.medical.weight} kg</div>
          <div>Blood Type: {person.medical.bloodType}</div>
          <div>
            <div className="font-medium">Medical Conditions</div>
            <div className="border rounded p-2 max-h-32 overflow-y-auto bg-[var(--card-bg)] whitespace-pre-line">
              {person.medical.medicalConditions || "None"}
            </div>
          </div>
        </>
      )}
    </div>
  )
}