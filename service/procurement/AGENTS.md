# BrewPipes "Procurement" Service Agent Guide

This guide is for agentic coding tools working in this particular service in the BrewPipes repo.
It captures commands and conventions observed in the current codebase.
When making changes to this service, its data models, its logic, or any other aspect, this document MUST be kept-up-to-date do that it accurately reflects the actual implementation.

## Service Domain

The Procurement service models supplier management and purchasing workflows: suppliers, purchase orders, line items, fees, shipping expectations, and receiving status.

## Overview

The big picture: the system tracks who you buy from and what you order, including quantities, costs, and receipt status.

Supplier
- A supplier is a vendor record with contact and address details.
- Suppliers are referenced by UUID in downstream services; purchase orders reference suppliers internally.

Purchase Order
- A purchase order is the top-level document for a supplier order.
- It includes a unique order number, status, expected arrival, and notes.

Purchase Order Line
- A line item captures what you are buying (ingredient, packaging, service, equipment, other).
- Each line includes quantity, unit, unit cost, currency, and a line number.
- Inventory items are referenced by opaque UUIDs for cross-service traceability.

Purchase Order Fee
- Fees capture non-line costs like shipping, tax, freight, or handling.
- Each fee has a type, amount, and currency.

Cross-domain references
- The Procurement service stores external identifiers (e.g., inventory item UUIDs) without foreign keys.
- This keeps the service deployable on its own database while preserving traceability.

## User Journey: Procurement Manager

Here is a simple procurement story that follows purchasing records, told in brewery terms.

You create a supplier for your hop vendor with contact and address details. You start a purchase order with an order number and the supplier, add line items for hops and yeast with quantities and unit costs, and include a shipping fee. While the order is being negotiated, it remains in draft. When you submit it to the supplier, the status updates to submitted, and once confirmed, it moves to confirmed with an expected arrival date.

As deliveries arrive, you update the order status to partially_received or received. The line items reference inventory item UUIDs so inventory can tie receipts back to the original order without sharing tables. Over time, the purchase order provides a clear paper trail of what was ordered, what it cost, and where it came from.

In short:
- suppliers define who you buy from,
- purchase orders track what you ordered,
- line items capture quantities and costs,
- fees cover shipping and other charges,
- statuses reflect ordering and receiving progress.

## Acceptance Criteria

- A procurement manager can create a supplier with a name and optional contact and address details.
- A purchase order can be created with a unique order number, supplier reference, and a default status of draft.
- Purchase orders support statuses for ordering and receiving flow (draft, submitted, confirmed, partially_received, received, cancelled).
- Line items can be added to a purchase order with item type, name, quantity, unit, unit cost, and currency.
- Line items can reference inventory items by opaque UUIDs without foreign keys.
- Fees can be attached to a purchase order with fee type, amount, and currency.
- Purchase orders can store expected arrival timestamps and notes for receiving coordination.
- Procurement records retain traceability for inventory and production without shared tables or cross-service foreign keys.
