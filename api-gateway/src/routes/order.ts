import { Router } from "express";
import axios from "axios";
import { config } from "../config";

const router = Router();

router.get("/", async (req, res) => {
  try {
    const response = await axios.get(`${config.orderServiceUrl}/orders`);
    res.json(response.data);
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      res.status(error.response.status).json(error.response.data);
    } else {
      res.status(500).json((error as Error).message);
    }
  }
});

router.post("/", async (req, res) => {
  try {
    const response = await axios.post(`${config.orderServiceUrl}/orders`, {
      userId: req.body.userId,
      status: req.body.status,
      items: req.body.items,
    });
    res.status(response.status).json(response.data);
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      res.status(error.response.status).json(error.response.data);
    } else {
      res.status(500).json((error as Error).message);
    }
  }
});

router.delete("/:orderId", async (req, res) => {
  try {
    await axios.delete(
      `${config.orderServiceUrl}/orders/${req.params.orderId}`
    );
    res.status(204).send();
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      res.status(error.response.status).json(error.response.data);
    } else {
      res.status(500).json((error as Error).message);
    }
  }
});

router.get("/:orderId", async (req, res) => {
  try {
    const response = await axios.get(
      `${config.orderServiceUrl}/orders/${req.params.orderId}`
    );
    res.json(response.data);
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      res.status(error.response.status).json(error.response.data);
    } else {
      res.status(500).json((error as Error).message);
    }
  }
});

router.patch("/:orderId/payment-status", async (req, res) => {
  try {
    await axios.patch(
      `${config.orderServiceUrl}/orders/${req.params.orderId}/payment-status`,
      null,
      {
        params: { status: req.query.status },
      }
    );
    res.status(204).send();
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      res.status(error.response.status).json(error.response.data);
    } else {
      res.status(500).json((error as Error).message);
    }
  }
});

export default router;
