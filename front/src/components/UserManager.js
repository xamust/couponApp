import React, { useState, useEffect, useCallback } from "react";
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
    Box,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    List,
    ListItem,
    ListItemText,
    Select,
    MenuItem,
    FormControl,
    InputLabel,
    TextField,
    CircularProgress,
} from "@mui/material";
import {
    createUser,
    getUsers,
    deleteUser,
    getCouponByUserId,
    applyCoupon,
    getCoupons,
} from "../api";

const UserManager = () => {
    const [users, setUsers] = useState([]);
    const [totalUsers, setTotalUsers] = useState(0);
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(5);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    // Состояние для модальных окон
    const [openViewCoupons, setOpenViewCoupons] = useState(false);
    const [openApplyCoupon, setOpenApplyCoupon] = useState(false);

    // Состояние для выбранного пользователя
    const [selectedUserCoupons, setSelectedUserCoupons] = useState([]);
    const [availableCoupons, setAvailableCoupons] = useState([]);
    const [selectedUserId, setSelectedUserId] = useState(null);
    const [selectedCouponId, setSelectedCouponId] = useState("");

    const [openCreateUser, setOpenCreateUser] = useState(false);
    const [newUser, setNewUser] = useState({ name: "", metadata: "" });

    const fetchUsers = useCallback(async (page, rowsPerPage) => {
        try {
            setLoading(true);
            const offset = page * rowsPerPage;
            const response = await getUsers(rowsPerPage, offset);
            setUsers(response.data || []);

            // костыль :(
            const resp2 = await getUsers(0, 0);
            setTotalUsers(resp2.data.length || 0);

            setError("");
        } catch (error) {
            handleError(error);
        } finally {
            setLoading(false);
        }
    }, []);

    const handleError = (error) => {
        if (error.response) {
            setError(error.response.data.message || "Ошибка сервера");
        } else if (error.request) {
            setError("Не удалось связаться с сервером");
        } else {
            setError("Произошла ошибка: " + error.message);
        }
    };

    const handleViewCoupons = async (userId) => {
        try {
            setSelectedUserId(userId);
            const response = await getCouponByUserId(userId);
            setSelectedUserCoupons(response.data || []);
            setOpenViewCoupons(true);
            setError("");
        } catch (error) {
            handleError(error);
        }
    };

    const handleOpenApplyCoupon = async (userId) => {
        try {
            setSelectedUserId(userId);
            const availableResponse = await getCoupons(0, 0);
            setAvailableCoupons(availableResponse.data || []);
            setOpenApplyCoupon(true);
            setError("");
        } catch (error) {
            handleError(error);
        }
    };

    const handleApplyCoupon = async () => {
        try {
            if (!selectedCouponId) {
                setError("Пожалуйста, выберите купон.");
                return;
            }
            await applyCoupon(selectedUserId, selectedCouponId);
            setSelectedCouponId("");
            setOpenApplyCoupon(false);
            setError("");
        } catch (error) {
            handleError(error);
        }
    };

    const handleDeleteUser = async (id) => {
        try {
            await deleteUser(id);
            fetchUsers(page, rowsPerPage);
            setError("");
        } catch (error) {
            handleError(error);
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

    const handleCreateUser = async () => {
        try {
            const metadata =
                newUser.metadata.trim() === "" ? {} : JSON.parse(newUser.metadata);
            await createUser({ name: newUser.name, metadata });
            setOpenCreateUser(false);
            setNewUser({ name: "", metadata: "" });
            setError("");
            fetchUsers(page, rowsPerPage);
        } catch (error) {
            if (error instanceof SyntaxError) {
                setError("Метаданные должны быть валидным JSON!");
            } else {
                handleError(error);
            }
        }
    };

    useEffect(() => {
        fetchUsers(page, rowsPerPage);
    }, [page, rowsPerPage, fetchUsers]);

    return (
        <Box>
            <Typography variant="h4" align="center" gutterBottom>
                Список пользователей
            </Typography>
            {error && (
                <Typography color="error" variant="body1" align="center" gutterBottom>
                    {error}
                </Typography>
            )}
            {loading ? (
                <Box display="flex" justifyContent="center" mt={4}>
                    <CircularProgress />
                </Box>
            ) : (
                <>
                    <TableContainer component={Paper}>
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell>Имя</TableCell>
                                    <TableCell>Активен</TableCell>
                                    <TableCell>Действия</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {users.length > 0 ? (
                                    users.map((user) => (
                                        <TableRow key={user.id}>
                                            <TableCell>{user.name}</TableCell>
                                            <TableCell>{user.is_active ? "Да" : "Нет"}</TableCell>
                                            <TableCell>
                                                <Button
                                                    color="primary"
                                                    size="small"
                                                    onClick={() => handleViewCoupons(user.id)}
                                                >
                                                    Посмотреть купоны
                                                </Button>
                                                <Button
                                                    color="secondary"
                                                    size="small"
                                                    onClick={() => handleOpenApplyCoupon(user.id)}
                                                >
                                                    Применить купон
                                                </Button>
                                                <Button
                                                    color="error"
                                                    size="small"
                                                    onClick={() => handleDeleteUser(user.id)}
                                                >
                                                    Удалить
                                                </Button>
                                            </TableCell>
                                        </TableRow>
                                    ))
                                ) : (
                                    <TableRow>
                                        <TableCell colSpan={3} align="center">
                                            Пользователи не найдены
                                        </TableCell>
                                    </TableRow>
                                )}
                            </TableBody>
                        </Table>
                    </TableContainer>
                    <TablePagination
                        component="div"
                        count={totalUsers}
                        page={page}
                        onPageChange={handleChangePage}
                        rowsPerPage={rowsPerPage}
                        onRowsPerPageChange={handleChangeRowsPerPage}
                        rowsPerPageOptions={[5, 10, 25, 50, 100]}
                    />
                </>
            )}
            <Button
                variant="contained"
                color="primary"
                onClick={() => setOpenCreateUser(true)}
                sx={{ mt: 2 }}
            >
                Создать пользователя
            </Button>

            <Dialog
                open={openCreateUser}
                onClose={() => setOpenCreateUser(false)}
                fullWidth
                maxWidth="sm"
            >
                <DialogTitle>Создать нового пользователя</DialogTitle>
                <DialogContent>
                    <TextField
                        label="Имя пользователя"
                        fullWidth
                        margin="normal"
                        value={newUser.name}
                        onChange={(e) =>
                            setNewUser({ ...newUser, name: e.target.value })
                        }
                    />
                    <TextField
                        label="Метаданные (JSON)"
                        fullWidth
                        margin="normal"
                        multiline
                        rows={4}
                        placeholder='Например, {"key": "value"}'
                        value={newUser.metadata}
                        onChange={(e) =>
                            setNewUser({ ...newUser, metadata: e.target.value })
                        }
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpenCreateUser(false)} color="primary">
                        Отмена
                    </Button>
                    <Button
                        onClick={handleCreateUser}
                        color="primary"
                        variant="contained"
                        disabled={!newUser.name.trim()}
                    >
                        Создать
                    </Button>
                </DialogActions>
            </Dialog>

            <Dialog
                open={openViewCoupons}
                onClose={() => setOpenViewCoupons(false)}
                fullWidth
                maxWidth="sm"
            >
                <DialogTitle>Применённые купоны</DialogTitle>
                <DialogContent>
                    {selectedUserCoupons.length > 0 ? (
                        <List>
                            {selectedUserCoupons.map((coupon) => (
                                <ListItem key={coupon.id}>
                                    <ListItemText
                                        primary={`${coupon.name} - ${coupon.reward}`}
                                    />
                                </ListItem>
                            ))}
                        </List>
                    ) : (
                        <Typography color="textSecondary">
                            Применённых купонов нет.
                        </Typography>
                    )}
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpenViewCoupons(false)} color="primary">
                        Закрыть
                    </Button>
                </DialogActions>
            </Dialog>

            <Dialog
                open={openApplyCoupon}
                onClose={() => setOpenApplyCoupon(false)}
                fullWidth
                maxWidth="sm"
            >
                <DialogTitle>Применить купон</DialogTitle>
                <DialogContent>
                    <FormControl fullWidth margin="normal">
                        <InputLabel id="coupon-select-label">Выберите купон</InputLabel>
                        <Select
                            labelId="coupon-select-label"
                            value={selectedCouponId}
                            onChange={(e) => setSelectedCouponId(e.target.value)}
                        >
                            {availableCoupons.map((coupon) => (
                                <MenuItem key={coupon.id} value={coupon.id}>
                                    {coupon.name}
                                </MenuItem>
                            ))}
                        </Select>
                    </FormControl>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpenApplyCoupon(false)} color="primary">
                        Отмена
                    </Button>
                    <Button
                        onClick={handleApplyCoupon}
                        color="primary"
                        variant="contained"
                        disabled={!selectedCouponId}
                    >
                        Применить
                    </Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
};

export default UserManager;
