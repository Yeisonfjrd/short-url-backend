import express from 'express';
import cors from 'cors';
import pkg from 'pg';
import crypto from 'crypto';

const { Pool } = pkg;
const app = express();
const port = process.env.PORT || 8080;

app.use(cors({
  origin: 'http://localhost:9000',
  credentials: true
}));
app.use(express.json());

const pool = new Pool({ 
  connectionString: process.env.DATABASE_URL 
});

app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

app.post('/api/shorten', async (req, res) => {
  try {
    const { url } = req.body;
    if (!url) {
      return res.status(400).json({ error: 'URL is required' });
    }

    const hash = crypto.createHash('sha256').update(url).digest('base64').slice(0, 8);

    const existingUrl = await pool.query(
      'SELECT * FROM urls WHERE original_url = $1',
      [url]
    );

    if (existingUrl.rows.length > 0) {
      const shortUrl = `${process.env.BASE_URL}/${existingUrl.rows[0].short_code}`;
      return res.json({ short_url: shortUrl });
    }

    const result = await pool.query(
      'INSERT INTO urls (original_url, short_code) VALUES ($1, $2) RETURNING *',
      [url, hash]
    );

    const shortUrl = `${process.env.BASE_URL}/${hash}`;
    res.json({ short_url: shortUrl });
  } catch (error) {
    console.error('Error creating short URL:', error);
    res.status(500).json({ error: 'Failed to create short URL' });
  }
});

app.get('/:shortCode', async (req, res) => {
  try {
    const { shortCode } = req.params;
    const result = await pool.query(
      'UPDATE urls SET visits = visits + 1, last_visit = NOW() WHERE short_code = $1 RETURNING original_url',
      [shortCode]
    );

    if (result.rows.length === 0) {
      return res.status(404).json({ error: 'URL not found' });
    }

    res.redirect(result.rows[0].original_url);
  } catch (error) {
    console.error('Error redirecting:', error);
    res.status(500).json({ error: 'Failed to redirect' });
  }
});

app.get('/api/stats/:shortCode', async (req, res) => {
  try {
    const { shortCode } = req.params;
    const result = await pool.query(
      'SELECT * FROM urls WHERE short_code = $1',
      [shortCode]
    );

    if (result.rows.length === 0) {
      return res.status(404).json({ error: 'URL not found' });
    }

    const url = result.rows[0];
    res.json({
      short_code: url.short_code,
      url: url.original_url,
      visits: url.visits,
      last_visit: url.last_visit,
      created_at: url.created_at
    });
  } catch (error) {
    console.error('Error getting stats:', error);
    res.status(500).json({ error: 'Failed to get stats' });
  }
});

app.listen(port, () => {
  console.log(`Server running on http://localhost:${port}`);
}); 