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
    private static long EXPIRATION_MS = 1000 * 60 * 60 * 24;//1 день
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
        try {
            Claims claims = Jwts.parser()
                    .verifyWith(SECRET_KEY)
                    .build()
                    .parseSignedClaims(token)
                    .getPayload();
//            System.out.println(claims.getExpiration());
            // Проверка срока действия
            if (claims.getExpiration().before(new Date())) {
                throw new ExpiredJwtException(null, claims, "Token expired");
            }

            return claims.getSubject();
        } catch (ExpiredJwtException ex) {
            System.err.println("Token expired.");
            // Обработка просроченного токена
            throw ex;
        } catch (JwtException | IllegalArgumentException e) {
            // Обработка других ошибок
            throw new JwtException("Invalid token");
        }
    }
}