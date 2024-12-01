import React, { useState, useEffect } from "react";
import {
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TablePagination,
    Button,
    Typography,
    Paper,
    TextField,
    Box,
} from "@mui/material";
import { getCoupons, deleteCoupon, createCoupon } from "../api";

const CouponManager = () => {
    const [coupons, setCoupons] = useState([]);
    const [totalCoupons, setTotalCoupons] = useState(0);
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(5);
    const [newCoupon, setNewCoupon] = useState({
        name: "",
        reward: "",
        maxRedemptions: 0,
        redeemBy: "",
    });
    const [error, setError] = useState("");

    const fetchCoupons = async (page, rowsPerPage) => {
        try {
            const offset = page * rowsPerPage;
            const response = await getCoupons(rowsPerPage, offset);
            const data = response.data;
            setCoupons(data || []);

            const resp2 = await getCoupons(0, 0);
            setTotalCoupons(resp2.data.length || 0);

            setError("");
        } catch (error) {
            if (error.response) {
                setError(error.response.data.message || "Ошибка сервера");
            } else if (error.request) {
                setError("Не удалось связаться с сервером");
            } else {
                setError("Произошла ошибка: " + error.message);
            }
        }
    };

    const handleCreateCoupon = async () => {
        try {
            const couponData = { ...newCoupon };
            if (couponData.maxRedemptions === 0 || !couponData.maxRedemptions) {
                delete couponData.maxRedemptions;
            }
            if (!couponData.redeemBy) {
                delete couponData.redeemBy;
            }
            await createCoupon(couponData);
            setNewCoupon({ name: "", reward: "", maxRedemptions: 0, redeemBy: "" });
            fetchCoupons(page, rowsPerPage);
            setError("");
        } catch (error) {
            if (error.response) {
                setError(error.response.data.message || "Ошибка сервера");
            } else if (error.request) {
                setError("Не удалось связаться с сервером");
            } else {
                setError("Произошла ошибка: " + error.message);
            }
        }
    };

    const handleDeleteCoupon = async (id) => {
        try {
            await deleteCoupon(id);
            fetchCoupons(page, rowsPerPage);
            setError("");
        } catch (error) {
            if (error.response) {
                setError(error.response.data.message || "Ошибка сервера");
            } else if (error.request) {
                setError("Не удалось связаться с сервером");
            } else {
                setError("Произошла ошибка: " + error.message);
            }
        }
    };

    const handleChangePage = (event, newPage) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (event) => {
        const newRowsPerPage = parseInt(event.target.value, 10);
        setRowsPerPage(newRowsPerPage);
        setPage(0);
    };

    useEffect(() => {
        fetchCoupons(page, rowsPerPage);
    }, [page, rowsPerPage]);

    return (
        <Box>
            <Typography variant="h4" align="center" gutterBottom>
                Список купонов
            </Typography>
            {error && (
                <Typography
                    color="error"
                    variant="body1"
                    align="center"
                    gutterBottom
                >
                    {error}
                </Typography>
            )}
            <TableContainer component={Paper}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>Название</TableCell>
                            <TableCell>Награда</TableCell>
                            <TableCell>Макс. Использований</TableCell>
                            <TableCell>Срок действия</TableCell>
                            <TableCell>Действия</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {coupons.length > 0 ? (
                            coupons.map((coupon) => (
                                <TableRow key={coupon.id}>
                                    <TableCell>{coupon.name}</TableCell>
                                    <TableCell>{coupon.reward}</TableCell>
                                    <TableCell>{coupon.maxRedemptions}</TableCell>
                                    <TableCell>
                                        {coupon.redeemBy
                                            ? new Date(coupon.redeemBy).toLocaleDateString()
                                            : "Не указано"}
                                    </TableCell>
                                    <TableCell>
                                        <Button
                                            color="error"
                                            onClick={() => handleDeleteCoupon(coupon.id)}
                                            size="small"
                                        >
                                            Удалить
                                        </Button>
                                    </TableCell>
                                </TableRow>
                            ))
                        ) : (
                            <TableRow>
                                <TableCell colSpan={5} align="center">
                                    Купоны не найдены
                                </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </TableContainer>
            <TablePagination
                component="div"
                count={totalCoupons}
                page={page}
                onPageChange={handleChangePage}
                rowsPerPage={rowsPerPage}
                onRowsPerPageChange={handleChangeRowsPerPage}
                rowsPerPageOptions={[5, 10, 25, 50, 100]}
            />
            <Box mt={4}>
                <Typography variant="h6">Добавить новый купон</Typography>
                <TextField
                    label="Название"
                    fullWidth
                    margin="normal"
                    value={newCoupon.name}
                    onChange={(e) =>
                        setNewCoupon({ ...newCoupon, name: e.target.value })
                    }
                />
                <TextField
                    label="Награда"
                    fullWidth
                    margin="normal"
                    value={newCoupon.reward}
                    onChange={(e) =>
                        setNewCoupon({ ...newCoupon, reward: e.target.value })
                    }
                />
                <TextField
                    label="Макс. Использований"
                    type="number"
                    fullWidth
                    margin="normal"
                    value={newCoupon.maxRedemptions}
                    onChange={(e) =>
                        setNewCoupon({
                            ...newCoupon,
                            maxRedemptions: parseInt(e.target.value, 10),
                        })
                    }
                />
                <TextField
                    label="Срок действия (YYYY-MM-DD)"
                    type="date"
                    fullWidth
                    margin="normal"
                    InputLabelProps={{ shrink: true }}
                    value={newCoupon.redeemBy}
                    onChange={(e) =>
                        setNewCoupon({ ...newCoupon, redeemBy: e.target.value })
                    }
                />
                <Button
                    variant="contained"
                    color="primary"
                    onClick={handleCreateCoupon}
                >
                    Создать купон
                </Button>
            </Box>
        </Box>
    );
};

export default CouponManager;
