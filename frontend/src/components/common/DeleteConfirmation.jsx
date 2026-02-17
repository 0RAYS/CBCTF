/**
 * Standard delete confirmation UI used in CRUD modals.
 * @param {Object} props
 * @param {string} props.message - Main confirmation message
 * @param {string} [props.warning] - Warning text shown below message in red
 * @param {string} [props.itemName] - Name of the item being deleted (shown bold)
 */
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
