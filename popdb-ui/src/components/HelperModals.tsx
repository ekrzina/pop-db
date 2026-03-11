"use client"

type BloodTypePickerModalProps = {
  value: string
  onSelect: (bloodType: string) => void
  onClose: () => void
}

export default function BloodTypePickerModal({ value, onSelect, onClose }: BloodTypePickerModalProps) {
  const bloodTypes = ["A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"]

  return (
    <div className="fixed inset-0 bg-black bg-opacity-30 flex justify-center items-center z-50">
      <div className="bg-white p-6 rounded shadow-lg space-y-4 w-80">
        <h3 className="text-lg font-semibold">Select Blood Type</h3>

        <div className="grid grid-cols-3 gap-2">
          {bloodTypes.map((bt) => (
            <button
              key={bt}
              className={`border px-3 py-2 rounded hover:bg-black hover:text-white ${value === bt ? "bg-black text-white" : ""}`}
              onClick={() => onSelect(bt)}
            >
              {bt}
            </button>
          ))}
        </div>

        <button
          className="mt-2 px-4 py-2 border rounded w-full"
          onClick={onClose}
        >
          Cancel
        </button>
      </div>
    </div>
  )
}