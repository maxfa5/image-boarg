package org.image_board.DTO;

import lombok.Data;

@Data
public class LoginResponseDTO {
  private String token;
  private String username;
  public LoginResponseDTO(String token, String username) {
    this.token = token;
    this.username = username;
  }
}

