package org.image_board.Utils;

import io.jsonwebtoken.*;
import io.jsonwebtoken.security.Keys;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;
import java.nio.charset.StandardCharsets;
import java.util.Date;
@Component
public class JwtUtils {
    private static SecretKey SECRET_KEY = null;
    private static final long EXPIRATION_MS = 864000000; // 10 дней
    public JwtUtils(@Value("${secret}") String secretKeyString) {
        // Преобразуем строку в ключ (для HS512 нужны 512 бит = 64 символа)
        SECRET_KEY = Keys.hmacShaKeyFor(secretKeyString.getBytes());
        System.out.println(SECRET_KEY);
    }

    public static String generateToken(String username) {
        if (SECRET_KEY == null) {
            System.err.println("Cannot generate token because SECRET_KEY is null.");
            return null; // Or throw an exception
        }
        Date now = new Date();
        Date expiryDate = new Date(now.getTime() + EXPIRATION_MS);
        return Jwts.builder()
                .subject(username)
                .issuedAt(now)
                .expiration(expiryDate)
                .signWith(SECRET_KEY)
                .compact();
    }

    public static String extractUsername(String token) {
        return Jwts.parser()
                .verifyWith(SECRET_KEY)
                .build()
                .parseSignedClaims(token)
                .getPayload()
                .getSubject();
    }
}