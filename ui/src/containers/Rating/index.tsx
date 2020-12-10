import React from 'react';
import { StarIcon } from '@patternfly/react-icons';
import './Rating.css';

const Rating: React.FC = () => {
  return (
    <div className="hub-details-rating">
      <form>
        <ul>
          {Array.apply(0, Array(5)).map((_, element) => (
            <li key={`${element + 1}-input`} className="hub-rate-area">
              <input type="radio" id={`${element}-star`} name="rating" value={`${element}`} />
              <label htmlFor={`${element}-star`}>
                <StarIcon />
              </label>
            </li>
          ))}
        </ul>
      </form>
    </div>
  );
};

export default Rating;
