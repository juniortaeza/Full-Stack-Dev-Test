import { useState } from "react";

import "../style/ExistingClient.css";

// custom formula: estimate range = [sum(bc) + sum(rate*min), sum(bc) + sum(rate*max)]
function getTotalEstimateCost(
  selectedEquipment,
  selectedLaborRate,
  propertyType,
) {
  let sumBaseCost = 0;
  for (const eqpmt of selectedEquipment) {
    sumBaseCost += eqpmt.baseCost;
  }

  let sumRateMin = 0.0;
  for (const rate of selectedLaborRate) {
    sumRateMin += rate.hourlyRate * rate.estimatedHours.min;
  }

  let sumRateMax = 0.0;
  for (const rate of selectedLaborRate) {
    sumRateMax += rate.hourlyRate * rate.estimatedHours.max;
  }

  if (propertyType == "residential") {
    // residential property cost no multiplier
    return [sumBaseCost + sumRateMin, sumBaseCost + sumRateMax];
  } else {
    // commercial seems to cost more, added arbitrary 1.5 multiplier
    return [(sumBaseCost + sumRateMin) * 1.5, (sumBaseCost + sumRateMax) * 1.5];
  }
}

function ExistingClient() {
  // const variables to capture data
  const [customerID, setCustomerID] = useState("");
  const [customer, setCustomer] = useState(null);

  const [equipment, setEquipment] = useState([]);
  const [selectedEquipment, setSelectedEquipment] = useState([]);

  const [laborRates, setLaborRates] = useState([]);
  const [selectedLaborRates, setSelectedLaborRates] = useState([]);

  // hit api endpoint returning equipment
  const loadEquipment = async () => {
    const response = await fetch("http://localhost:8080/api/equipment");
    const data = await response.json();
    setEquipment(data);
  };

  // hit api endpoint returning labor rates
  const loadLaborRates = async () => {
    const response = await fetch("http://localhost:8080/api/laborRates");
    const data = await response.json();
    setLaborRates(data);
  };

  // capture user-selected equipment and store
  const handleEquipmentChange = (item) => {
    setSelectedEquipment((current) => {
      if (current.includes(item)) {
        return current.filter((equipment) => equipment !== item);
      }

      return [...current, item];
    });
  };

  // capture user-selected rates and store
  const handleLaborChange = (rate) => {
    setSelectedLaborRates((current) => {
      if (current.includes(rate)) {
        return current.filter((laborRate) => laborRate !== rate);
      }

      return [...current, rate];
    });
  };

  // match general keywords to categories, this is currently hardcoded
  const getSystemCategories = (systemType) => {
    const keywords = {
      "Central AC": "Air Conditioner",
      "Gas Furnace": "Furnace",
      "Gas Heater": "Furnace",
      "Heat Pump": "Heat Pump",
      "Mini-Split": "Mini-Split",
      Rooftop: "Rooftop Unit",
      "Package Unit": "Package Unit",
    };

    const categories = [];

    for (const keyword in keywords) {
      if (systemType.includes(keyword)) {
        categories.push(keywords[keyword]);
      }
    }

    return categories;
  };

  // return list of relevant equipment based off system type
  const relevantEquipment = customer
    ? equipment.filter((item) =>
        getSystemCategories(customer.systemType).includes(item.category),
      )
    : [];

  // search for customer in database, populate equipment and rate lists
  const searchCustomer = async () => {
    const response = await fetch(
      `http://localhost:8080/api/customers/${customerID}`,
    );

    const data = await response.json();

    setSelectedEquipment([]);
    setSelectedLaborRates([]);

    setCustomer(data);
    await loadEquipment();
    await loadLaborRates();
  };

  return (
    <div>
      <div>
        <h3>Enter Client ID</h3>
        <input
          type="text"
          placeholder="ex: CUST001"
          value={customerID}
          onChange={(e) => setCustomerID(e.target.value)}
        />

        <button onClick={searchCustomer}>Search</button>
      </div>

      <div>
        {customer !== null && (
          <div className="info-card">
            <p>
              <b>Name</b>: {customer.name}
            </p>
            <p>
              <b>Address</b>: {customer.address}
            </p>
            <p>
              <b>Phone</b>: {customer.phone}
            </p>
            <p>
              <b>Property Type</b>: {customer.propertyType}
            </p>
            <p>
              <b>System Type</b>: {customer.systemType}
            </p>
          </div>
        )}

        <div className="checklist-row">
          {relevantEquipment.length > 0 && (
            <div className="info-card option-card">
              <h3>Suggested Equipment</h3>
              <div className="equipment-list">
                {relevantEquipment.map((item) => (
                  <label className="equipment-item" key={item.id}>
                    <input
                      type="checkbox"
                      checked={selectedEquipment.includes(item)}
                      onChange={() => handleEquipmentChange(item)}
                    />
                    <span className="equipment-name">{item.name}</span>
                    <span className="equipment-price">${item.baseCost}</span>
                  </label>
                ))}
              </div>
            </div>
          )}

          {laborRates.length > 0 && (
            <div className="info-card option-card">
              <h3>Labor Rates</h3>

              <div className="labor-list">
                {laborRates.map((rate) => (
                  <label
                    className="labor-item"
                    key={`${rate.jobType}-${rate.level}`}
                  >
                    <input
                      type="checkbox"
                      checked={selectedLaborRates.includes(rate)}
                      onChange={() => handleLaborChange(rate)}
                    />

                    <span className="labor-name">
                      {rate.jobType} - {rate.level}
                    </span>

                    <span className="labor-price">${rate.hourlyRate}/hr</span>
                  </label>
                ))}
              </div>
            </div>
          )}

          {selectedEquipment.length > 0 && selectedLaborRates.length > 0 && (
            <div className="info-card option-card">
              <h3>Labor Rates</h3>
              <h2>
                Property Type: <u>{customer.propertyType}</u>
              </h2>
              <h1>
                $
                {getTotalEstimateCost(
                  selectedEquipment,
                  selectedLaborRates,
                  customer.propertyType,
                )[0].toFixed(2)}{" "}
                - $
                {getTotalEstimateCost(
                  selectedEquipment,
                  selectedLaborRates,
                  customer.propertyType,
                )[1].toFixed(2)}
              </h1>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default ExistingClient;
