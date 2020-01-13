const XFormModal = () => {
    const [isOpen, setIsOpen] = useState(false);
    const closeModal = () => {
        setIsOpen(false);
    }
    const openModal = () => {
        setIsOpen(true);
    }

    const [params, setParams] = useState({ name: '', content: '', tags: [], targets: [] });
    const handleChange = (e, { name, value }) => {
        console.log("Updated form values: ", name, value, params);
        setParams({ ...params, [name]: value });
    }

    return (
        <Modal
            centered={false}
            open={isOpen}
            onClose={closeModal}
            trigger={<Button positive circular icon='plus' onClick={openModal} />}

            // Form properties
            as={Form}
            onSubmit={handleSubmit}
        ></Modal>
    );
}