/** Standard delete confirmation UI used in CRUD modals. */
function DeleteConfirmation({ message, warning, itemName }) {
  return (
    <div className="text-center py-4">
      <p className="text-neutral-50">
        {message}
        {itemName && <span className="font-bold"> {itemName}</span>}
      </p>
      {warning && <p className="text-red-400 text-sm mt-2">{warning}</p>}
    </div>
  );
}

export default DeleteConfirmation;
