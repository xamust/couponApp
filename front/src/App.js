import React from "react";
import { Container, Grid, Typography } from "@mui/material";
import CouponManager from "./components/CouponManager";
import UserManager from "./components/UserManager";

function App() {
    return (
        <Container maxWidth="lg">
            <Typography variant="h4" align="center" gutterBottom>
                Административная панель - Управление купонами
            </Typography>
                <Grid item xs={12} md={6}>
                    <UserManager />
                </Grid>
                <Grid item xs={12} md={6}>
                    <CouponManager />
                </Grid>
        </Container>
    );
}

export default App;