import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'http://localhost:8080/api';

export let options = {
    vus: 10, // n виртуальных пользователей
    duration: '5s', // Тест идет t секунд
};

export default function () {
    let username = `optoed_${__VU}`;
    let password = 'optoed';

    // 0. Регистрация/Авторизация пользователя
    let authPayload = JSON.stringify({ username: username, password: password });
    let authHeaders = { 'Content-Type': 'application/json' };
    let authRes = http.post(`${BASE_URL}/auth`, authPayload, { headers: authHeaders });

    check(authRes, {
        'Регистрация/Авторизация успешна': (res) => res.status === 200,
    });

    if (authRes.status !== 200) return;

    // Логируем тело ответа, чтобы понять, что пошло не так
    //console.log('authRes:', authRes.body);

    let token = JSON.parse(authRes.body).token;
    let authHeadersWithToken = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
    };

    // 1. Покупка товара
    let buyRes = http.post(`${BASE_URL}/buy/pen`, {}, { headers: authHeadersWithToken });

    check(buyRes, {
        'Покупка успешна': (res) => res.status === 200 || res.status === 400, //не хватило денег
    });

    // Логируем тело ответа, чтобы понять, что пошло не так
    //console.log('Buy response:', buyRes.body);

    // 2. Получение информации о пользователе
    let infoRes = http.get(`${BASE_URL}/info`, { headers: authHeadersWithToken });

    check(infoRes, {
        'Информация получена': (res) => res.status === 200,
    });

    //let coins = JSON.parse(infoRes.body).schema.coins;

    // Логируем тело ответа, чтобы понять, что пошло не так
    //console.log('infoRes:', infoRes.body);

    // 3. Передача монет другому пользователю
    let sendCoinPayload = JSON.stringify({ toUser: "optoed_1", amount: 1 });
    let sendCoinRes = http.post(`${BASE_URL}/sendCoin`, sendCoinPayload, { headers: authHeadersWithToken });

    check(sendCoinRes, {
        'Передача монет успешна': (res) => res.status === 200 || res.status === 400, // 400 если недостаточно монет
    });

    // Логируем тело ответа, чтобы понять, что пошло не так
    //console.log('sendCoinRes:', sendCoinRes.body);

    sleep(1);
}
