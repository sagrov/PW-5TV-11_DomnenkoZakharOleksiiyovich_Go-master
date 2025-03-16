document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('form');
    const result = document.getElementById('result');

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        let inputs;
        const url = form.id === 'calculator1' ? '/api/calculator1' : '/api/calculator2';
        console.log(url)
        if (form.id === 'calculator1') {
            // Convert all inputs to strings (trimmed)
            inputs = [...form.querySelectorAll('input')].map(input => input.value);
        } else if (form.id === 'calculator2') {
            // Convert all inputs to floats, ignoring non-numeric values
            inputs = [...form.querySelectorAll('input')].map(input => {
                const value = parseFloat(input.value);
                return isNaN(value) ? 0 : value; // Replace NaN with 0 or handle as needed
            });
        }

        try {
            const response = await fetch(url, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ values: inputs })
            });

            const data = await response.json();
            console.log(data);
            result.textContent = `Result: ${data.result}`;
        } catch (error) {
            result.textContent = 'Error calculating.';
            console.error(error);
        }
    });
});
