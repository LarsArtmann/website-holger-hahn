#!/bin/bash

# Deployment script for Holger Hahn Website to Google Cloud Run
# Usage: ./deploy.sh [PROJECT_ID] [REGION]

set -e

# Configuration
PROJECT_ID=${1:-"your-project-id"}
REGION=${2:-"us-central1"}
SERVICE_NAME="holger-hahn-website"
IMAGE_NAME="gcr.io/$PROJECT_ID/$SERVICE_NAME"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Deploying Holger Hahn Website to Google Cloud Run${NC}"
echo -e "${BLUE}Project: $PROJECT_ID${NC}"
echo -e "${BLUE}Region: $REGION${NC}"
echo ""

# Check if gcloud is installed and authenticated
if ! command -v gcloud &> /dev/null; then
    echo -e "${RED}❌ gcloud CLI is not installed${NC}"
    echo "Please install it from: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Check if user is authenticated
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    echo -e "${YELLOW}⚠️  Not authenticated with gcloud${NC}"
    echo "Please run: gcloud auth login"
    exit 1
fi

# Set the project
echo -e "${YELLOW}📋 Setting gcloud project to $PROJECT_ID${NC}"
gcloud config set project $PROJECT_ID

# Enable required APIs
echo -e "${YELLOW}🔧 Enabling required APIs${NC}"
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com

# Build and deploy using Cloud Build
echo -e "${YELLOW}🏗️  Building and deploying with Cloud Build${NC}"
gcloud builds submit --config cloudbuild.yaml \
    --substitutions=_REGION=$REGION

# Get the service URL
echo -e "${YELLOW}🔍 Getting service URL${NC}"
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region=$REGION --format="value(status.url)")

echo ""
echo -e "${GREEN}✅ Deployment completed successfully!${NC}"
echo -e "${GREEN}🌐 Service URL: $SERVICE_URL${NC}"
echo ""
echo -e "${BLUE}📊 Useful commands:${NC}"
echo -e "${BLUE}  View logs: gcloud run services logs tail $SERVICE_NAME --region=$REGION${NC}"
echo -e "${BLUE}  Get status: gcloud run services describe $SERVICE_NAME --region=$REGION${NC}"
echo -e "${BLUE}  Update traffic: gcloud run services update-traffic $SERVICE_NAME --to-latest --region=$REGION${NC}"
echo ""

# Optional: Test the deployment
read -p "🧪 Would you like to test the deployment? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}🧪 Testing deployment...${NC}"
    if curl -f -s "$SERVICE_URL/health" > /dev/null; then
        echo -e "${GREEN}✅ Health check passed!${NC}"
    else
        echo -e "${RED}❌ Health check failed${NC}"
    fi
    
    echo -e "${YELLOW}📄 Fetching homepage...${NC}"
    if curl -f -s "$SERVICE_URL" > /dev/null; then
        echo -e "${GREEN}✅ Homepage accessible!${NC}"
    else
        echo -e "${RED}❌ Homepage not accessible${NC}"
    fi
fi

echo -e "${GREEN}🎉 Deployment process complete!${NC}"