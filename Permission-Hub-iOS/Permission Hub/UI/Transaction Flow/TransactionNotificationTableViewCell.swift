//
//  TransactionNotificationTableViewCell.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

enum TransactionNotificationType {
    case verification, warning

    var image: UIImage? {
        switch self {
        case .verification:
            return UIImage(named: "checkmark_small")
        case .warning:
            return UIImage(named: "warning")
        }
    }

    var color: UIColor {
        switch self {
        case .verification:
            return PHColors.topaz
        case .warning:
            return PHColors.red
        }
    }
}

final class TransactionNotificationTableViewCell: UITableViewCell {

    // MARK: - Private properties

    private lazy var stackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .horizontal
        stackView.spacing = 10
        stackView.alignment = .center

        return stackView
    }()

    private lazy var notificationImageView: UIImageView = {

        let imageView = UIImageView()
        imageView.translatesAutoresizingMaskIntoConstraints = false
        imageView.setContentHuggingPriority(.required, for: .horizontal)

        return imageView
    }()

    private lazy var disclosureImageView: UIImageView = {

        let image = UIImage(named: "disclosure_indicator")
        let imageView = UIImageView(image: image)
        imageView.translatesAutoresizingMaskIntoConstraints = false
        imageView.setContentHuggingPriority(.required, for: .horizontal)

        return imageView
    }()

    private lazy var notificationLabel: UILabel = {

        let label = UILabel()
        label.font = PHFonts.wesBold()
        label.textColor = .white
        label.textAlignment = .left

        return label
    }()

    private lazy var notificationSublabel: UILabel = {

        let label = UILabel()
        label.font = PHFonts.wesRegular(ofSize: 10)
        label.textColor = .white

        return label
    }()

    private lazy var separatorView: UIView = {

        let view = UIView()
        view.translatesAutoresizingMaskIntoConstraints = false

        view.backgroundColor = .white

        return view
    }()

    // MARK: - Initialization

    override init(style: UITableViewCell.CellStyle, reuseIdentifier: String?) {
        super.init(style: style, reuseIdentifier: reuseIdentifier)

        configure()
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Configuration

    private func configure() {

        selectionStyle = .none

        addSubview(stackView)
        addSubview(separatorView)

        let margin: CGFloat = 15
        stackView.topAnchor.constraint(equalTo: topAnchor, constant: margin).isActive = true
        stackView.leftAnchor.constraint(equalTo: leftAnchor, constant: margin).isActive = true
        stackView.rightAnchor.constraint(equalTo: rightAnchor, constant: -margin).isActive = true

        separatorView.topAnchor.constraint(equalTo: stackView.bottomAnchor, constant: margin).isActive = true
        separatorView.leftAnchor.constraint(equalTo: leftAnchor).isActive = true
        separatorView.rightAnchor.constraint(equalTo: rightAnchor).isActive = true
        separatorView.bottomAnchor.constraint(equalTo: bottomAnchor).isActive = true
        separatorView.heightAnchor.constraint(equalToConstant: 1).isActive = true

        stackView.addArrangedSubview(notificationImageView)
        stackView.addArrangedSubview(notificationLabel)
        stackView.addArrangedSubview(notificationSublabel)
        stackView.addArrangedSubview(disclosureImageView)
    }

    func configure(withType type: TransactionNotificationType, andText text: String) {

        notificationImageView.image = type.image
        backgroundColor = type.color
        notificationLabel.text = text
        notificationSublabel.text = "See details"
    }
}
