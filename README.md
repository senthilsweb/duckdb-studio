# duckdb-studio

DuckDB Studio! A simple, yet powerful web utility designed to explore and interact with DuckDB databases. This tool is perfect for developers, data analysts, and hobbyists who seek an easy and efficient way to work with DuckDB databases either locally or remotely. DuckDB Studio serves as an extension and complementary solution to the [DuckDB Data API](https://github.com/senthilsweb/duckdb_data_api), enhancing its capabilities and user experience. 

## Features

DuckDB Studio offers a range of features designed to simplify your data exploration:
- **Easy Viewing:** Access and view databases and tables seamlessly.
- **Table Exploration:** Navigate through your data with a friendly table view and pagination.
- **SQL Query Execution:** Execute any SQL query directly from the interface.
- **Prettify SQL:** Beautify your SQL queries for better readability and maintenance.
- **Versatile Deployment:** Host DuckDB Studio on Vercel for free or run it on any computer with NodeJS.

## Getting Started

To get started with DuckDB Studio, you'll need to set up your development environment. Here are the environment variables required to configure your instance:

```.env
NUXT_SESSION_PASSWORD=your_secure_password
NUXT_PUBLIC_LOGO_WEB=path/to/your/web/logo.svg
NUXT_PUBLIC_LOGO_MOBILE=path/to/your/mobile/logo.svg
NUXT_PUBLIC_DUCKDB_DATA_API_BASE_PATH=your_duckdb_data_api_endpoint
```

### Development Environment

1. **Clone the repository:**

```bash
git clone https://github.com/senthilsweb/duckdb-studio
cd duckdb-studio
```

2. **Install dependencies:**

```bash
npm install
```

3. **Run the development server:**

```bash
npm run dev
```

Visit `http://localhost:3000` to see DuckDB Studio in action.

### Production Deployment

DuckDB Studio can be deployed on various platforms. A common choice is Vercel, which offers free hosting:

1. **Build your project for production:**

```bash
npm run build
```

2. **Start the production server:**

```bash
npm start
```

3. **Deploy to Vercel:**

- Push your changes to a GitHub repository.
- Connect your GitHub repository to Vercel.
- Follow Vercel's prompts to deploy your DuckDB Studio instance.

## To-Do Items

- [ ] Automatic online API documentation for tables/entities.
- [ ] Integrate SQLGlot for query transpilation to multiple data sources.
- [ ] Develop an SQL Query Optimizer.
- [ ] Implement features for enhanced data lineage insights.
- [ ] Go binary for packing the app as appliance



DuckDB Studio is inspired by and aims to complement the DuckDB ecosystem by providing an accessible and user-friendly interface for data exploration and management. Whether you're a seasoned data professional or just starting out, DuckDB Studio is designed to streamline your workflow and enhance your data interaction experiences.

