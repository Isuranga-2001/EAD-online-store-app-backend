import dotenv from "dotenv";

dotenv.config();

export const config = {
  productServiceUrl: process.env.PRODUCT_SERVICE_URL as string,
  orderServiceUrl: process.env.ORDER_SERVICE_URL as string,
  userServiceUrl: process.env.USER_SERVICE_URL as string,
  paymentServiceUrl: process.env.PAYMENT_SERVICE_URL as string,
  fileServiceUrl: process.env.FILE_SERVICE_URL as string,
  cartServiceUrl: process.env.CART_SERVICE_URL as string,
};
