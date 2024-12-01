import axios from "axios";

const api = axios.create({
    baseURL: "/api/v1", // Прокси настроена
});

// Получение списка купонов
export const getCoupons = (limit,offset) => api.post("/coupon", { limit: limit, offset: offset });

// Создание нового купона
export const createCoupon = (couponData) => api.post("/coupon/create", couponData);

// Удаление купона
export const deleteCoupon = (id) => api.delete(`/coupon/${id}`);

// применение купона
export const applyCoupon = (user_id,coupon_id) => api.post(`/coupon/apply/`,{user_id:user_id,coupon_id:coupon_id});

// получение купона по user_id
export const getCouponByUserId = (user_id) => api.get(`/coupon/apply/${user_id}`);

// Получение списка пользователей
export const getUsers = (limit,offset) => api.post("/user", { limit: limit, offset: offset });

// Создание нового пользователя
export const createUser = (userData) => api.post("/user/create", userData);

// Удаление пользователя
export const deleteUser = (id) => api.delete(`/user/${id}`);


export default api;