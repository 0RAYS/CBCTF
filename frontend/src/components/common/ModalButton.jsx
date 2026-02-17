import Button from './Button';

const ModalButton = ({
  children,
  onClick,
  variant = 'default', // 'default' | 'primary' | 'danger'
  disabled = false,
}) => {
  const variants = {
    default: 'ghost',
    primary: 'primary',
    danger: 'danger',
  };

  return (
    <Button variant={variants[variant]} size="sm" onClick={onClick} disabled={disabled}>
      {children}
    </Button>
  );
};

export default ModalButton;
