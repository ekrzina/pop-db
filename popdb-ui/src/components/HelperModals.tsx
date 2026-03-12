"use client"

type BloodTypePickerModalProps = {
  value: string
  onSelect: (bloodType: string) => void
  onClose: () => void
}

export default function BloodTypePickerModal({ value, onSelect, onClose }: BloodTypePickerModalProps) {
  const bloodTypes = ["A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"]

  return (
    <div className="fixed inset-0 bg-black/40 backdrop-blur-sm flex justify-center items-center p-4 z-50">
      <div className="bg-white rounded-2xl shadow-xl w-full max-w-sm flex flex-col">

        {/* Header */}
        <div className="flex justify-between items-center px-5 py-3 border-b">
          <h3 className="text-lg font-semibold">Select Blood Type</h3>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 text-xl leading-none"
          >
            ×
          </button>
        </div>

        {/* Options */}
        <div className="p-5 grid grid-cols-3 gap-3">
          {bloodTypes.map((bt) => (
            <button
              key={bt}
              onClick={() => onSelect(bt)}
              className={`
                px-3 py-2 rounded-lg border text-center transition
                ${value === bt
                  ? "bg-rose-500 text-white border-rose-500"
                  : "bg-gray-50 hover:bg-gray-100"}
              `}
            >
              {bt}
            </button>
          ))}
        </div>

        {/* Footer */}
        <div className="px-5 py-3 border-t">
          <button
            onClick={onClose}
            className="w-full px-4 py-2 border rounded-lg hover:bg-gray-100 transition"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  )
}
