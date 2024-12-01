import React from "react";
import { Container, Grid, Typography } from "@mui/material";
import CouponManager from "./components/CouponManager";
import UserManager from "./components/UserManager";

function App() {
    return (
        <Container>
            <Typography variant="h4" align="center" gutterBottom>
                Административная панель - Управление купонами
            </Typography>
            <Grid container spacing={2}>
                <Grid item xs={6}>
                    <UserManager />
                </Grid>
                <Grid item xs={6}>
                    <CouponManager />
                </Grid>
            </Grid>
        </Container>
    );
}

export default App;