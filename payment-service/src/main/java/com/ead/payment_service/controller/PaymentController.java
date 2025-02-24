package com.ead.payment_service.controller;

import com.ead.payment_service.dto.PaymentCreateDTO;
import com.ead.payment_service.dto.PaymentUpdateDTO;
import com.ead.payment_service.dto.PaymentDTO;
import com.ead.payment_service.service.PaymentService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import jakarta.validation.Valid;
import java.util.List;

@RestController
@RequestMapping("/payments")
public class PaymentController {

    @Autowired
    private PaymentService paymentService;

    @PostMapping
    public ResponseEntity<PaymentDTO> createPayment(@RequestBody @Valid PaymentCreateDTO paymentCreateDTO) {
        PaymentDTO payment = paymentService.createPayment(paymentCreateDTO);
        return ResponseEntity.status(HttpStatus.CREATED).body(payment);
    }

    @GetMapping("/{id}")
    public ResponseEntity<PaymentDTO> getPaymentById(@PathVariable Long id) {
        PaymentDTO payment = paymentService.getPaymentById(id);
        return ResponseEntity.ok(payment);
    }

    @GetMapping("/order/{orderId}")
    public ResponseEntity<List<PaymentDTO>> getPaymentsByOrderId(@PathVariable Long orderId) {
        List<PaymentDTO> payments = paymentService.getPaymentsByOrderId(orderId);
        return ResponseEntity.ok(payments);
    }

    @PutMapping("/{id}")
    public ResponseEntity<PaymentDTO> updatePayment(@PathVariable Long id, @RequestBody @Valid PaymentUpdateDTO paymentUpdateDTO) {
        paymentUpdateDTO.setPaymentId(id);
        PaymentDTO updatedPayment = paymentService.updatePayment(paymentUpdateDTO);
        return ResponseEntity.ok(updatedPayment);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<String> deletePayment(@PathVariable Long id) {
        paymentService.deletePaymentById(id);
        return ResponseEntity.ok("Payment deleted successfully.");
    }

    @DeleteMapping("/order/{orderId}")
    public ResponseEntity<String> deletePaymentsByOrderId(@PathVariable Long orderId) {
        paymentService.deletePaymentsByOrderId(orderId);
        return ResponseEntity.ok("Payments deleted successfully.");
    }
}