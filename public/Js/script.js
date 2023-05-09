// fungsi button ketika diklik mengarah ke email, di halaman contact

function getData() {
    event.preventDefault()
    let nama = document.getElementById('name').value;
    let email = document.getElementById('email').value
    let phone = document.getElementById('nohandphone').value
    let subject = document.getElementById('subject').value
    let message = document.getElementById('message').value

    if(nama == '') {
        alert`SILAHKAN ISIKAN NAMA ANDA TERLEBIH DAHULU`
    } else if (email == '') {
        alert`JANGAN LUPA EMAILNYA, SUPAYA KITA BISA TERHUBUNG`
    } else if (message == '') {
        alert`PESAN NYA JANGAN LUPA YA!`
    }

    let data = {
        nama,
        email,
        phone,
        subject,
        message
    };
    
    const myEmail = 'ibnualinsani23@gmail.com'
    let a = document.createElement('a');
    a.setAttribute('href', `https://mail.google.com/mail/?view=cm&fs=1&to=${myEmail}&su=${data.subject}&body=Halo, nama saya ${data.nama}. ${data.message}
    hubungi saya di nomor ini: ${data.phone}`);
    a.click();
}


