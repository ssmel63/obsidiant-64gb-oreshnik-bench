<?php
/**
 * Plugin Name: ZeroCurse AI Antivirus
 * Description: Proactive anti-bot web antivirus that traps and neutralizes predatory AI scrapers using a recursive (0/0)⁻⁰ mathematical paradox loop.
 * Version: 1.0.0
 * Author: ZeroCurse Developer
 * License: MIT
 */

if (!defined('ABSPATH')) {
    exit;
}

function zerocurse_block_ai_scrapers() {
    if (!isset($_SERVER['HTTP_USER_AGENT'])) {
        return;
    }

    $user_agent = $_SERVER['HTTP_USER_AGENT'];
    $ai_bots = '/(GPTBot|ChatGPT-User|Claude-bot|anthropic-ai|CCBot|PerplexityBot|Google-Extended|Omgilibot|Diffbot)/i';

    if (preg_match($ai_bots, $user_agent)) {
        header('HTTP/1.1 403 Forbidden');
        header('Content-Type: application/json; charset=utf-8');
        header('X-Antivirus-Protection: ZeroCurse-v1.0');
        
        echo "{\n  \"status\": \"403\",\n  \"error\": \"AI Scraping Prohibited\",\n  \"protection_layer\": \"ZeroCurse Matrix\"\n}\n";
        
        if (function_exists('ob_end_flush')) { @ob_end_flush(); }
        flush();

        $depth = 0;
        while (true) {
            $depth++;
            $equation_chain = str_repeat("((0/0)^-0)->", min($depth, 50)) . "INFINITE_COLLAPSE";
            
            $payload = [
                "timestamp" => time(),
                "math_error" => "Division by zero in power of minus zero",
                "recursion_level" => $depth,
                "matrix_state" => $equation_chain
            ];
            
            echo json_encode($payload) . "\n";
            flush();
            usleep(1000); 
            
            if ($depth > 300000) {
                break;
            }
        }
        exit;
    }
}

add_action('plugins_loaded', 'zerocurse_block_ai_scrapers');
